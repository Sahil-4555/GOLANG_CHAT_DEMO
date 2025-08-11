import PropTypes from 'prop-types';
import { useCallback, useContext, useEffect, useRef, useState } from 'react';
import api from '../utils/api';
import { connect, useDispatch } from 'react-redux';
import * as constants from '../utils/constant/Constant';
import WebSocketContext from '../utils/context/WebSocketContext';
import "../assets/style/components/GroupUserPopup.css"

function GroupUserPopup(props) {
    const {
        searchedvalue,
        setSearchedvalue,
        onClose,
        groupUserPopup,
        setAddeduser,
        addeduser,
        selectedChannelState,
        usersState

    } = props;
    const userId = localStorage.getItem("userId");
    const dispatch = useDispatch();
    const [focusedIndex, setFocusedIndex] = useState(0);
    const [isGroupAdmin, setGroupAdmin] = useState(false);
    const [listing, setListing] = useState(true);
    const [preventUsers, setPreventUsers] = useState([]);
    const groupUserListRef = useRef(null);
    const groupUserAddRef = useRef(null);
    const groupUserContainerRef = useRef(null);
    const groupListContainerRef = useRef(null);
    const { sendMessage, lastMessage, readyState } = useContext(WebSocketContext);

    const fetchPreventUsers = useCallback(async () => {
        if (!selectedChannelState?._id) return [];
        try {
            const result = await api.get(`v1/chat/get-channel-members/${selectedChannelState?._id}`);
            if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
                const data = result?.data?.data[0]?.users || [];
                setPreventUsers(data);
                for (let i = 0; i < data.length; i++) {
                    if (data[i].user_id === userId && data[i].is_admin) {
                        setGroupAdmin(true)
                        break;
                    }
                }
                return data;
            }
        } catch (err) {
            console.log('Error:', err);
            return [];
        }
    }, [selectedChannelState?._id]);

    const filterUser = usersState?.filter(
        (ele) =>
            !addeduser.includes(ele) && !preventUsers.some((preventUser) => preventUser.user_id === ele._id)
    );

    useEffect(() => {
        fetchPreventUsers();
    }, [fetchPreventUsers, selectedChannelState?._id]);

    const listFocus = (e) => {
        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                setFocusedIndex((prevIndex) => (prevIndex < filterUser.length - 1 ? prevIndex + 1 : 0));
                break;
            case 'ArrowUp':
                e.preventDefault();
                setFocusedIndex((prevIndex) => (prevIndex > 0 ? prevIndex - 1 : filterUser.length - 1));
                break;
            case 'Enter':
                if (focusedIndex !== -1) {
                    handleSelect(filterUser[focusedIndex]);
                }
                break;
            case 'Backspace':
                if (searchedvalue === '' && addeduser.length > 0) {
                    setAddeduser((prevAddedUsers) => {
                        const updatedUsers = [...prevAddedUsers];
                        updatedUsers.pop();
                        return updatedUsers;
                    });
                }
                break;
            default:
                break;
        }
    };

    const handleKeyDown = (e) => {
        if (e.key === 'Escape') {
            onClose(false);
        }
    };

    const handleSelect = (option) => {
        if (!addeduser.some((user) => user._id === option._id)) {
            setAddeduser((prevAddedUsers) => [...prevAddedUsers, option]);
        }
        setSearchedvalue('');
        setFocusedIndex(0);
    };

    const handleAddAdmin = async (channelId, userId) => {
        try {
            const requestBody = {
                channel_id: channelId,
                user_id: userId
            };
            const result = await api.put(`v1/channel/give-admin-rights`, requestBody)
            if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
                const data = result?.data?.data[0]?.users || [];
                setPreventUsers(data);
            }
        } catch (err) {
            console.log('Error: ', err)
        }
    }

    useEffect(() => {
        if (focusedIndex !== -1 && groupUserAddRef.current && groupUserContainerRef.current) {
            const listItem = groupUserAddRef.current.children[focusedIndex];
            if (listItem) {
                listItem.scrollIntoView({
                    block: 'nearest',
                    inline: 'nearest',
                    behavior: 'smooth'
                });
            }
        }
    }, [focusedIndex]);

    const handleDeselect = (option) => {
        setAddeduser((prevAddedUsers) => prevAddedUsers.filter((user) => user._id !== option._id));
    };

    const handleClickOutside = (e) => {
        const insideClickElement = document.getElementById('insideClick');
        const isOutsideContainers =
            (!groupUserContainerRef.current || !groupUserContainerRef.current.contains(e.target)) &&
            (!groupListContainerRef.current || !groupListContainerRef.current.contains(e.target));

        if (isOutsideContainers && (!insideClickElement || !insideClickElement.contains(e.target))) {
            onClose(false);
        }
    };

    useEffect(() => {
        if (!groupUserPopup) return;
        queueMicrotask(() => {
            groupUserListRef?.current?.focus();
            setSearchedvalue('');
            setFocusedIndex(0);
        });
    }, [groupUserPopup, setSearchedvalue]);

    useEffect(() => {
        const detectClickElement = document.getElementById('detectclick');
        document.addEventListener('keydown', handleKeyDown);
        detectClickElement.addEventListener('click', handleClickOutside);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
            detectClickElement.removeEventListener('click', handleClickOutside);
        };
    }, []);

    const updateChannel = async () => {
        if (!addeduser.length) return;
        try {

            // Event to send to the added members to reflect the channel on channel list
            const channelId = selectedChannelState?._id;
            const memberIds = addeduser.map((user) => user._id);
            const WebSocketPayload = {
                action: constants.WEBSOCKET_ACTIONS.POST_GLOBAL_NOTIFICATION_TO_USERS,
                notification_type: constants.WEBSOCKET_ACTIONS.ADD_CHANNEL_ON_ADDING_MEMBER,
                target: channelId,
                payload: selectedChannelState,
                notify_users: memberIds,
            }
            var data = JSON.stringify(WebSocketPayload)
            sendMessage(data)

            // Add added members in the database 
            const requestBody = {
                channel_id: channelId,
                users: memberIds
            };
            await api.put('v1/channel/addmembers', requestBody);

            const userIdsFromSelectedChannel = selectedChannelState.users.map((user) => user.user_id);
            const channelMembers = [...memberIds, ...userIdsFromSelectedChannel];
            const WebSocketData = {
                action: constants.WEBSOCKET_ACTIONS.POST_GLOBAL_NOTIFICATION_TO_USERS,
                notification_type: constants.WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_USERS,
                target: channelId,
                payload: addeduser,
                notify_users: channelMembers,
            };
            data = JSON.stringify(WebSocketData)
            sendMessage(data)


            onClose();
            fetchPreventUsers();
        } catch (error) {
            console.error('Error updating channel:', error.message);
        }
    };

    return (
        <div className={`popup-container ${groupUserPopup ? 'visible' : 'hidden'}`} id="detectclick">
            <div className='popup-content' id='insideClick'>
                {listing ? (
                    <div className='pop-content-list'>
                        <div className='popup-header'>
                            <h5 className='popup-channelname'>Channels members - {selectedChannelState?.channel_name}</h5>
                            <button onClick={onClose}>
                                <svg
                                    fill='#000000'
                                    height='20px'
                                    width='20px'
                                    version='1.1'
                                    id='Capa_1'
                                    xmlns='http:/www.w3.org/2000/svg'
                                    viewBox='0 0 490 490'
                                    xmlSpace='preserve'>
                                    <g id='SVGRepo_bgCarrier' strokeWidth='0'></g>
                                    <g
                                        id='SVGRepo_tracerCarrier'
                                        strokeLinecap='round'
                                        strokeLinejoin='round'
                                        stroke='#CCCCCC'
                                        strokeWidth='0.9800000000000001'></g>
                                    <g id='SVGRepo_iconCarrier'>
                                        {' '}
                                        <polygon points='456.851,0 245,212.564 33.149,0 0.708,32.337 212.669,245.004 0.708,457.678 33.149,490 245,277.443 456.851,490 489.292,457.678 277.331,245.004 489.292,32.337 '></polygon>{' '}
                                    </g>
                                </svg>
                            </button>
                        </div>
                        <div className='channel-user-container' ref={groupListContainerRef}>
                            {preventUsers.length > 0 ? (
                                <ul
                                    role='list'
                                    className='channel-users-list'>
                                    {preventUsers.map((person) => {
                                        return (
                                            <li
                                                key={person.user_id}
                                                className='channel-user-list'>
                                                <div className='channel-user-avatar '>
                                                    <div className='channel-user-avatar-name'>
                                                        {person.name.charAt(0).toUpperCase()}
                                                    </div>
                                                    <span
                                                        className={`channel-user-status`}>
                                                        {person?.status === constants.USER_STATUS.ONLINE && (
                                                            <div className='userStatusOnlineUserList'></div>
                                                        )}
                                                        {person?.status === constants.USER_STATUS.AWAY && (
                                                            <div className='userStatusAwayUserList'></div>
                                                        )}
                                                        {person?.status === constants.USER_STATUS.DO_NOT_DISTURB && (
                                                            <div className='userStatusDonotDisturbUserList'></div>
                                                        )}
                                                        {person?.status === constants.USER_STATUS.OFFLINE && (
                                                            <div className='userStatusOfflineUserList'></div>
                                                        )}
                                                    </span>

                                                    <div className='channel-user-details'>
                                                        <p className='channel-user-details-name'>
                                                            {person.name} {person.is_admin && <span className='channel-user-admin-tag'>(Admin)</span>}
                                                        </p>
                                                        <p className='channel-user-details-email'>
                                                            {person.email}
                                                        </p>
                                                    </div>
                                                </div>
                                                {
                                                    isGroupAdmin && !person.is_admin && (
                                                        <button
                                                            onClick={() => handleAddAdmin(selectedChannelState?._id, person.user_id)}
                                                            className='add-admin-button'>
                                                            Add as Admin
                                                        </button>
                                                    )
                                                }
                                            </li>
                                        );
                                    })}
                                </ul>
                            ) : (
                                <div className='channel-user-not-found'>No User Found</div>
                            )}
                            <div className='channel-user-buttons'>
                                <button onClick={onClose} className='channel-user-cancel-button'>
                                    Cancel
                                </button>
                                {isGroupAdmin && (
                                    <button
                                        onClick={() => setListing(false)}
                                        className='channel-user-add-button'>
                                        + Add
                                    </button>
                                )}
                            </div>
                        </div>
                    </div>
                ) : (
                    <div>
                        <div className='popup-header'>
                            <h5 className='popup-channelname'>Add Users - {selectedChannelState?.channel_name}</h5>
                            <button onClick={onClose}>
                                <svg
                                    fill='#000000'
                                    height='20px'
                                    width='20px'
                                    version='1.1'
                                    id='Capa_1'
                                    xmlns='http://www.w3.org/2000/svg'
                                    viewBox='0 0 490 490'
                                    xmlSpace='preserve'>
                                    <g id='SVGRepo_bgCarrier' strokeWidth='0'></g>
                                    <g
                                        id='SVGRepo_tracerCarrier'
                                        strokeLinecap='round'
                                        strokeLinejoin='round'
                                        stroke='#CCCCCC'
                                        strokeWidth='0.9800000000000001'></g>
                                    <g id='SVGRepo_iconCarrier'>
                                        {' '}
                                        <polygon points='456.851,0 245,212.564 33.149,0 0.708,32.337 212.669,245.004 0.708,457.678 33.149,490 245,277.443 456.851,490 489.292,457.678 277.331,245.004 489.292,32.337 '></polygon>{' '}
                                    </g>
                                </svg>
                            </button>
                        </div>

                        <div
                            className={`selected-options selectedOptions ${addeduser.length > 0 && 'p-2'} `}>
                            {addeduser.map((option, index) => (
                                <div
                                    key={index}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        handleDeselect(option);
                                    }}
                                    className='selected-option selectedOption'>
                                    {option?.name} <span className='remove-option'>✕</span>
                                </div>
                            ))}
                            <input
                                type='search'
                                ref={groupUserListRef}
                                value={searchedvalue}
                                onKeyDown={(e) => listFocus(e)}
                                onChange={(e) => setSearchedvalue(e.target.value)}
                                className={`search-user ${addeduser.length > 0 ? 'w-auto' : 'w-full'}`}
                                placeholder={
                                    addeduser.length > 0
                                        ? ''
                                        : addeduser.length === 0
                                            ? 'Please enter name'
                                            : 'Select option'
                                }
                            />
                        </div>

                        <div className='group-add-user-list-container' ref={groupUserContainerRef}>
                            {filterUser.length > 0 ? (
                                <ul
                                    role='list'
                                    className='filtered-user-list'
                                    ref={groupUserAddRef}>
                                    {filterUser?.map((users, index) => {
                                        return (
                                            <li
                                                key={index}
                                                onClick={() => {
                                                    handleSelect(users);
                                                }}
                                                className={`focused-user ${focusedIndex === index ? 'bg-blue-100' : ''}`}>
                                                <div className='channel-user-avatar'>
                                                    <div className='channel-user-avatar-name'>
                                                        {users?.name?.charAt(0).toUpperCase()}
                                                    </div>
                                                    <span
                                                        className={`channel-user-status`}>
                                                        {users?.status === constants.USER_STATUS.ONLINE && (
                                                            <div className='userStatusOnlineUserList'></div>
                                                        )}
                                                        {users?.status === constants.USER_STATUS.AWAY && (
                                                            <div className='userStatusAwayUserList'></div>
                                                        )}
                                                        {users?.status === constants.USER_STATUS.DO_NOT_DISTURB && (
                                                            <div className='userStatusDonotDisturbUserList'></div>
                                                        )}
                                                        {users?.status === constants.USER_STATUS.OFFLINE && (
                                                            <div className='userStatusOfflineUserList'></div>
                                                        )}
                                                    </span>
                                                    <div className='channel-user-details'>
                                                        <p className='channel-user-details-name'>
                                                            {users.name}
                                                        </p>
                                                        <p className='channel-user-details-email'>
                                                            {users.email}
                                                        </p>
                                                    </div>
                                                </div>
                                            </li>
                                        );
                                    })}
                                </ul>
                            ) : (
                                <div className='add-user-not-found'>
                                    No Result Found
                                </div>
                            )}
                        </div>

                        <div className='channel-user-buttons'>
                            <button onClick={onClose} className='channel-user-cancel-button'>
                                Cancel
                            </button>
                            <button
                                onClick={updateChannel}
                                disabled={!addeduser.length > 0}
                                className='add-user-save-button'>
                                Save
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </div >
    );
}

GroupUserPopup.propTypes = {
    searchedvalue: PropTypes.any.isRequired,
    setSearchedvalue: PropTypes.func.isRequired,
    onClose: PropTypes.func.isRequired,
    groupUserPopup: PropTypes.any.isRequired,
    setAddeduser: PropTypes.func.isRequired,
    addeduser: PropTypes.any.isRequired,
};

const mapStateToProps = state => ({
    selectedChannelState: state.SelectedChannelReducer || {},
    usersState: state.UsersReducer || [],
});

export default connect(
    mapStateToProps
)(GroupUserPopup);

