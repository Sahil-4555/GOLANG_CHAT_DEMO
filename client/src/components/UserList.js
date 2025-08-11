import { Menu, Transition } from '@headlessui/react';
import PropTypes from 'prop-types';
import React, { Fragment, useState } from 'react';
import {useDispatch, connect } from "react-redux";
import * as constants from '../utils/constant/Constant';
import cross from '../assets/icons/cross.svg'
import star from '../assets/icons/star.svg'
import threedots from '../assets/icons/threedots.svg'
import logouticon from '../assets/icons/logout.svg'
import * as handlers from '../redux/actions/actions'
import "../assets/style/components/UserList.css";

const UserList = ({
    onSelectUser,
    handleLeaveChannel,
    OpenAddUserPopup,
    isOpen,
    onLogout,
    searchOpen,
    handleChangeModal,
    directChannelsState,
    selectedChannelState,
    groupChannelsState,
    favouriteChannelsState
}) => {
    const userName = localStorage.getItem('userName');
    const handleLogout = () => {
        onLogout();
    };

    const [groupChannelDropdown, setGroupChannelDropdown] = useState(false);
    const [favDropdown, setFavdropdown] = useState(false);
    const [directChannelDropdown, setDirectChannelDropdown] = useState(false);
    const dispatch = useDispatch();

    return (
        <div className={`user-list-container ${isOpen ? 'translate-x-0 w-full' : '-translate-x-full w-full'}`}>
            {/* Header */}
            <div className='divison-border'>
                <div className='username-container'>
                    <p className='username-modal'>{userName}</p>
                    <div className='add-channel-container' onClick={() => handleChangeModal()}>
                        <p className='add-icon'>+</p>
                        <div className='add-channel-modal'>
                            Add Channel
                        </div>
                    </div>
                </div>
            </div>
            <div className='user-list-modal'>
                <div className='user-list-modal-wrapper'>
                    {/* Favourites Channels */}
                    {favouriteChannelsState.length > 0 && (
                        <div className='pt-1'>
                            <div
                                className='channel-dropdown'
                                onClick={() => setFavdropdown(!favDropdown)}>
                                <div className={`${favDropdown ? '-rotate-180' : 'rotate-0'} duration-300 `}>
                                    <svg
                                        xmlns='http://www.w3.org/2000/svg'
                                        fill='none'
                                        viewBox='0 0 20 20'
                                        strokeWidth={1.5}
                                        stroke='white'
                                        className='w-4 h-4'>
                                        <path
                                            strokeLinecap='round'
                                            strokeLinejoin='round'
                                            d='M19.5 8.25l-7.5 7.5-7.5-7.5'
                                        />
                                    </svg>
                                </div>
                                <h2 className='channel-title'>Favorites</h2>
                            </div>
                            <ul className={`${favDropdown ? 'my-2' : 'h-full'} duration-400 `}>
                                {favouriteChannelsState.map((fav) => {
                                    const isChannel = selectedChannelState?._id === fav._id;

                                    if (fav.channel_type === constants.CHANNEL_TYPES.DIRECT_CHANNEL) {
                                        return (
                                            <li
                                                key={fav?._id}
                                                onClick={() => { onSelectUser(fav, constants.SELECTION_TYPES.FROM_LIST); }}
                                                className={`channel-list-container group
                                                ${favDropdown && 'hide-channels'}
                                                ${isChannel ? 'selected-channel' : 'hover-over-channel'}`}>
                                                <div
                                                    className={`channel-avatar-one-to-one group
                                                    ${isChannel ? 'text-white' : 'hover:text-black'}`}>
                                                    {fav?.channel_name.charAt(0).toUpperCase()}
                                                    <span
                                                        className={`user-status-icon`}>
                                                        {fav?.users[0]?.status === constants.USER_STATUS.ONLINE && (
                                                            <div className='userStatusOnlineUserList'></div>
                                                        )}
                                                        {fav?.users[0]?.status === constants.USER_STATUS.AWAY && (
                                                            <div className='userStatusAwayUserList'></div>
                                                        )}
                                                        {fav?.users[0]?.status === constants.USER_STATUS.DO_NOT_DISTURB && (
                                                            <div className='userStatusDonotDisturbUserList'></div>
                                                        )}
                                                        {fav?.users[0]?.status === constants.USER_STATUS.OFFLINE && (
                                                            <div className='userStatusOfflineUserList'></div>
                                                        )}
                                                    </span>
                                                </div>
                                                <div className='channel-list-channelname'>
                                                    <span
                                                        className={`ml-2 ${isChannel ? 'text-white' : ' text-gray-400 group-hover:text-black'
                                                            }`}>
                                                        {fav?.channel_name.toUpperCase()}
                                                    </span>

                                                    {fav?.message_count > 0 && (
                                                        <span className='channel-message-count'>
                                                            {fav?.message_count}
                                                        </span>
                                                    )}

                                                    <span onClick={(e) => e?.stopPropagation()}>
                                                        <Menu as='div' className='menu-container'>
                                                            <div>
                                                                <Menu.Button className='menu-icon-three-dots'>
                                                                    <img src={threedots} alt="" style={{ width: 25, height: 25 }} />
                                                                </Menu.Button>
                                                            </div>
                                                            <Transition
                                                                as={Fragment}
                                                                enter='transition ease-out duration-100'
                                                                enterFrom='transform opacity-0 scale-95'
                                                                enterTo='transform opacity-100 scale-100'
                                                                leave='transition ease-in duration-75'
                                                                leaveFrom='transform opacity-100 scale-100'
                                                                leaveTo='transform opacity-0 scale-95'>
                                                                <Menu.Items className='menu-items'>
                                                                    <div className='px-1 py-1 '>
                                                                        <Menu.Item>
                                                                            {({ active }) => (
                                                                                <button
                                                                                    className={`${active ? 'activated' : 'deactivated'
                                                                                        } group active-menu`}
                                                                                    onClick={() => dispatch(handlers.handleRemovePrivateFavourites(fav))}>
                                                                                    <img src={star} alt="" style={{ width: 21, height: 21 }} />
                                                                                    <span>Unfavorite</span>
                                                                                </button>
                                                                            )}
                                                                        </Menu.Item>
                                                                        <Menu.Item>
                                                                            {({ active }) => (
                                                                                <button
                                                                                    onClick={() => dispatch(handlers.handleCloseConversation(fav, constants.SECTIONS.FAVOURITE_CHANNEL_SECTION, selectedChannelState?._id))}
                                                                                    className={`${active ? 'activated' : 'deactivated'
                                                                                        } group active-menu`}>
                                                                                    <img src={cross} alt="" style={{ width: 21, height: 21 }} />
                                                                                    <span>Close Conversation</span>
                                                                                </button>
                                                                            )}
                                                                        </Menu.Item>
                                                                    </div>
                                                                </Menu.Items>
                                                            </Transition>
                                                        </Menu>
                                                    </span>
                                                </div>
                                            </li>
                                        );
                                    } else {
                                        return (
                                            <li
                                                key={fav?._id}
                                                onClick={() => {
                                                    !isChannel && onSelectUser(fav, constants.SELECTION_TYPES.FROM_LIST);
                                                }}
                                                className={`channel-list-group-container group
                                                ${favDropdown && 'hide-channels'}
                                                ${isChannel ? 'selected-channel' : 'hover-over-channel'}`}>
                                                <div className='channel-avatar-group-channel'>
                                                    <svg
                                                        xmlns='http://www.w3.org/2000/svg'
                                                        fill='none'
                                                        viewBox='0 0 24 24'
                                                        stroke='currentColor'
                                                        className={`w-8 h-8 ml-0.5 ${isChannel ? 'text-white' : 'text-gray-400 group-hover:text-black'
                                                            }`}>
                                                        <path
                                                            strokeLinecap='round'
                                                            strokeLinejoin='round'
                                                            d='M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z'
                                                        />
                                                    </svg>
                                                    <div className='channel-list-channelname'>
                                                        <div
                                                            className={`ml-2 w-max ${isChannel ? 'text-white' : 'text-gray-400 group-hover:text-black'
                                                                }`}>
                                                            {fav?.channel_name.toUpperCase()}{' '}
                                                        </div>
                                                        {fav?.message_count > 0 && (
                                                            <span className='channel-message-count'>
                                                                {fav?.message_count}
                                                            </span>
                                                        )}
                                                        <span onClick={(e) => e?.stopPropagation()}>
                                                            <Menu as='div' className='menu-container'>
                                                                <div>
                                                                    <Menu.Button className='menu-icon-three-dots'>
                                                                        <img src={threedots} alt="" style={{ width: 25, height: 25 }} />
                                                                    </Menu.Button>
                                                                </div>
                                                                <Transition
                                                                    as={Fragment}
                                                                    enter='transition ease-out duration-100'
                                                                    enterFrom='transform opacity-0 scale-95'
                                                                    enterTo='transform opacity-100 scale-100'
                                                                    leave='transition ease-in duration-75'
                                                                    leaveFrom='transform opacity-100 scale-100'
                                                                    leaveTo='transform opacity-0 scale-95'>
                                                                    <Menu.Items className='menu-items'>
                                                                        <div className='px-1 py-1 '>
                                                                            <Menu.Item>
                                                                                {({ active }) => (
                                                                                    <button
                                                                                        className={`${active ? 'activated' : 'deactivated'
                                                                                            } group active-menu`}
                                                                                        onClick={() => dispatch(handlers.handleRemoveGroupFavourites(fav))}>
                                                                                        <img src={star} alt="" style={{ width: 21, height: 21 }} />
                                                                                        <span>Unfavorite</span>
                                                                                    </button>
                                                                                )}
                                                                            </Menu.Item>
                                                                            <Menu.Item>
                                                                                {({ active }) => (
                                                                                    <button
                                                                                        onClick={() => handleLeaveChannel(fav, constants.SECTIONS.FAVOURITE_CHANNEL_SECTION)}
                                                                                        className={`${active ? 'activated' : 'deactivated'
                                                                                            } group active-menu`}>
                                                                                        <img src={cross} alt="" style={{ width: 21, height: 21 }} />
                                                                                        <span>Leave Channel</span>
                                                                                    </button>
                                                                                )}
                                                                            </Menu.Item>
                                                                        </div>
                                                                    </Menu.Items>
                                                                </Transition>
                                                            </Menu>
                                                        </span>
                                                    </div>
                                                </div>
                                            </li>
                                        );
                                    }
                                })}
                            </ul>
                        </div>
                    )}

                    {/* Group Channels */}
                    <div className='pt-5'>
                        <div
                            className='channel-dropdown'
                            onClick={() => setGroupChannelDropdown(!groupChannelDropdown)}>
                            <div className={`${groupChannelDropdown ? '-rotate-180' : 'rotate-0'} duration-300 `}>
                                <svg
                                    xmlns='http://www.w3.org/2000/svg'
                                    fill='none'
                                    viewBox='0 0 20 20'
                                    strokeWidth={1.5}
                                    stroke='white'
                                    className='w-4 h-4'>
                                    <path
                                        strokeLinecap='round'
                                        strokeLinejoin='round'
                                        d='M19.5 8.25l-7.5 7.5-7.5-7.5'
                                    />
                                </svg>
                            </div>
                            <h2 className='channel-title'>Channels</h2>
                        </div>
                        <ul className={`${groupChannelDropdown ? 'my-2' : 'h-full '} duration-400 `}>
                            {groupChannelsState.map((channel) => {
                                const isChannel = selectedChannelState?._id === channel._id;
                                return (
                                    <li
                                        key={channel?._id}
                                        onClick={() => {
                                            !isChannel && onSelectUser(channel, constants.SELECTION_TYPES.FROM_LIST);
                                        }}
                                        className={`channel-list-group-container group
                                        ${groupChannelDropdown && 'hide-channels'}
                                        ${isChannel ? 'selected-channel' : 'hover-over-channel'}`}>
                                        <div className='channel-avatar-group-channel'>
                                            <svg
                                                xmlns='http://www.w3.org/2000/svg'
                                                fill='none'
                                                viewBox='0 0 24 24'
                                                stroke='currentColor'
                                                className={`w-8 h-8 ml-0.5 ${isChannel ? 'text-white' : 'text-gray-400 group-hover:text-black'
                                                    }`}>
                                                <path
                                                    strokeLinecap='round'
                                                    strokeLinejoin='round'
                                                    d='M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z'
                                                />
                                            </svg>
                                            <div className='channel-list-channelname '>
                                                <div
                                                    className={`ml-2 w-max ${isChannel ? 'text-white' : 'text-gray-400 group-hover:text-black'
                                                        }`}>
                                                    {channel?.channel_name.toUpperCase()}{' '}
                                                </div>
                                                {channel?.message_count > 0 && (
                                                    <span className='channel-message-count'>
                                                        {channel?.message_count}
                                                    </span>
                                                )}
                                                <span onClick={(e) => e?.stopPropagation()}>
                                                    <Menu as='div' className='menu-container'>
                                                        <div>
                                                            <Menu.Button className='menu-icon-three-dots'>
                                                                <img src={threedots} alt="" style={{ width: 25, height: 25 }} />
                                                            </Menu.Button>
                                                        </div>
                                                        <Transition
                                                            as={Fragment}
                                                            enter='transition ease-out duration-100'
                                                            enterFrom='transform opacity-0 scale-95'
                                                            enterTo='transform opacity-100 scale-100'
                                                            leave='transition ease-in duration-75'
                                                            leaveFrom='transform opacity-100 scale-100'
                                                            leaveTo='transform opacity-0 scale-95'>
                                                            <Menu.Items className='menu-items'>
                                                                <div className='px-1 py-1 '>
                                                                    <Menu.Item>
                                                                        {({ active }) => (
                                                                            <button
                                                                                className={`${active ? 'activated' : 'deactivated'
                                                                                    } group active-menu`}
                                                                                onClick={() => dispatch(handlers.handleGroupAddFavourites(channel))}>
                                                                                <img src={star} alt="" style={{ width: 21, height: 21 }} />
                                                                                <span>Favorite</span>
                                                                            </button>
                                                                        )}
                                                                    </Menu.Item>
                                                                    <Menu.Item>
                                                                        {({ active }) => (
                                                                            <button
                                                                                onClick={() => handleLeaveChannel(channel, constants.SECTIONS.GROUP_CHANNEL_SECTION)}
                                                                                className={`${active ? 'activated' : 'deactivated'
                                                                                    } group active-menu`}>
                                                                                <img src={cross} alt="" style={{ width: 21, height: 21 }} />
                                                                                <span>Leave Channel</span>
                                                                            </button>
                                                                        )}
                                                                    </Menu.Item>
                                                                </div>
                                                            </Menu.Items>
                                                        </Transition>
                                                    </Menu>
                                                </span>
                                            </div>
                                        </div>
                                    </li>
                                );
                            })}
                        </ul>
                    </div>

                    {/* Direct Channels */}
                    <div className='pt-5'>
                        <div className='direct-message-modal'>
                            <div
                                className='channel-dropdown'
                                onClick={() => setDirectChannelDropdown(!directChannelDropdown)}>
                                <div className={`${directChannelDropdown ? '-rotate-180' : 'rotate-0'} duration-300`}>
                                    <svg
                                        xmlns='http://www.w3.org/2000/svg'
                                        fill='none'
                                        viewBox='0 0 20 20'
                                        strokeWidth={1.5}
                                        stroke='white'
                                        className='w-4 h-4'>
                                        <path
                                            strokeLinecap='round'
                                            strokeLinejoin='round'
                                            d='M19.5 8.25l-7.5 7.5-7.5-7.5'
                                        />
                                    </svg>
                                </div>
                                <h2 className='channel-title'>
                                    Direct Messages
                                </h2>
                            </div>
                            <button onClick={searchOpen} className='ml-1'>
                                <svg
                                    xmlns='http://www.w3.org/2000/svg'
                                    fill='none'
                                    viewBox='0 0 24 24'
                                    strokeWidth={1.5}
                                    stroke='white'
                                    className='w-6 h-6'>
                                    <path strokeLinecap='round' strokeLinejoin='round' d='M12 4.5v15m7.5-7.5h-15' />
                                </svg>
                            </button>
                        </div>
                        <ul className={`${directChannelDropdown ? 'my-2' : 'h-full '} duration-400 `}>
                            {directChannelsState?.length > 0 &&
                                directChannelsState?.map((user) => {
                                    const isUser = selectedChannelState?._id === user._id;

                                    return (
                                        <li
                                            key={user?._id}
                                            onClick={() => { onSelectUser(user, constants.SELECTION_TYPES.FROM_LIST); }}
                                            className={`channel-list-container group
                                                ${directChannelDropdown && 'hide-channels'}
                                                ${isUser ? 'selected-channel' : 'hover-over-channel'}`}>
                                            <div
                                                className={`channel-avatar-one-to-one group
                                                    ${isUser ? 'text-white' : 'hover:text-black'}`}>
                                                {user?.channel_name.charAt(0).toUpperCase()}
                                                <span
                                                    className={`user-status-icon`}>
                                                    {user?.users[0]?.status === constants.USER_STATUS.ONLINE && (
                                                        <div className='userStatusOnlineUserList'></div>
                                                    )}
                                                    {user?.users[0]?.status === constants.USER_STATUS.AWAY && (
                                                        <div className='userStatusAwayUserList'></div>
                                                    )}
                                                    {user?.users[0]?.status === constants.USER_STATUS.DO_NOT_DISTURB && (
                                                        <div className='userStatusDonotDisturbUserList'></div>
                                                    )}
                                                    {user?.users[0]?.status === constants.USER_STATUS.OFFLINE && (
                                                        <div className='userStatusOfflineUserList'></div>
                                                    )}
                                                </span>
                                            </div>
                                            <div className='channel-list-channelname'>
                                                <span
                                                    className={`ml-2 ${isUser ? 'text-white' : ' text-gray-400 group-hover:text-black'
                                                        }`}>
                                                    {user?.channel_name.toUpperCase()}
                                                </span>

                                                {user?.message_count > 0 && (
                                                    <span className='channel-message-count'>
                                                        {user?.message_count}
                                                    </span>
                                                )}

                                                <span onClick={(e) => e?.stopPropagation()}>
                                                    <Menu as='div' className='menu-container'>
                                                        <div>
                                                            <Menu.Button className='menu-icon-three-dots'>
                                                                <img src={threedots} alt="" style={{ width: 25, height: 25 }} />
                                                            </Menu.Button>
                                                        </div>
                                                        <Transition
                                                            as={Fragment}
                                                            enter='transition ease-out duration-100'
                                                            enterFrom='transform opacity-0 scale-95'
                                                            enterTo='transform opacity-100 scale-100'
                                                            leave='transition ease-in duration-75'
                                                            leaveFrom='transform opacity-100 scale-100'
                                                            leaveTo='transform opacity-0 scale-95'>
                                                            <Menu.Items className='menu-items'>
                                                                <div className='px-1 py-1 '>
                                                                    <Menu.Item>
                                                                        {({ active }) => (
                                                                            <button
                                                                                className={`${active ? 'activated' : 'deactivated'
                                                                                    } group active-menu`}
                                                                                onClick={() => dispatch(handlers.handlePrivateAddFavourites(user))}>
                                                                                <img src={star} alt="" style={{ width: 21, height: 21 }} />
                                                                                <span>Favourite</span>
                                                                            </button>
                                                                        )}
                                                                    </Menu.Item>
                                                                    <Menu.Item>
                                                                        {({ active }) => (
                                                                            <button
                                                                                onClick={() => dispatch(handlers.handleCloseConversation(user, constants.SECTIONS.DIRECT_CHANNEL_SECTION, selectedChannelState?._id))}
                                                                                className={`${active ? 'activated' : 'deactivated'
                                                                                    } group active-menu`}>
                                                                                <img src={cross} alt="" style={{ width: 21, height: 21 }} />
                                                                                <span>Close Conversation</span>
                                                                            </button>
                                                                        )}
                                                                    </Menu.Item>
                                                                </div>
                                                            </Menu.Items>
                                                        </Transition>
                                                    </Menu>
                                                </span>
                                            </div>
                                        </li>
                                    );
                                })}
                        </ul>
                    </div>
                </div>
            </div>

            {/* Logout Button */}
            <div className='logout-button-container '>
                <div className='logout-button-wrapper' onClick={handleLogout}>
                    <img src={logouticon} alt="logout icon" style={{ width: 25, height: 25 }} />
                    <div className='logout-button-text'>
                        Logout
                    </div>
                </div>
            </div>
        </div>
    );
};

UserList.propTypes = {
    onSelectUser: PropTypes.func.isRequired,
    handleLeaveChannel: PropTypes.func.isRequired,
    isOpen: PropTypes.bool.isRequired,
    onLogout: PropTypes.func.isRequired,
    searchOpen: PropTypes.func.isRequired,
    OpenAddUserPopup: PropTypes.func.isRequired,
    handleChangeModal: PropTypes.func.isRequired,
};

const mapStateToProps = state => ({
    directChannelsState: state.DirectChannelReducer || [],
    groupChannelsState: state.GroupChannelReducer || [],
    favouriteChannelsState: state.FavouriteChannelReducer || [],
    selectedChannelState: state.SelectedChannelReducer || {},
});

export default connect(
    mapStateToProps
)(UserList);

