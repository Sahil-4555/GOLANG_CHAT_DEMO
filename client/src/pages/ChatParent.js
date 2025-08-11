// ChatParent.jsx
import React, { useEffect, useState, useContext, useCallback } from "react";
import { useDispatch } from "react-redux";
import { ReadyState } from "react-use-websocket";
import * as handlers from "../redux/actions/actions";
import * as types from "../utils/constant/ActionTypes"
import UserList from "../components/UserList";
import ChannelModal from "../components/ChannelModal";
import WebSocketContext from "../utils/context/WebSocketContext";
import { API_STATUS, CHANNEL_TYPES, SECTIONS, SELECTION_TYPES, USER_STATUS, WEBSOCKET_ACTIONS } from "../utils/constant/Constant";
import api from "../utils/api";
import SearchUser from "../components/SearchUser";
import ChatScreen from "./ChatScreen";
import GroupUserPopup from "../components/GroupUserPopup";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { connect } from 'react-redux';
import Loading from "../components/Loading";
import "../index.css"
import ConfirmationModal from "../components/ActionPopup";

const ChatParent = ({ selectedChannelState }) => {
  const toastOptions = {
    position: "top-right",
    autoClose: 3000,
    pauseOnHover: true,
    draggable: true,
    theme: "dark",
  };

  const userId = localStorage.getItem("userId");
  const dispatch = useDispatch();
  const [searchedChannel, setSearchedChannel] = useState([]);
  const [page, setPage] = useState(1);
  const [searchuservalue, setSearchuservalue] = useState('');
  const [searchgroupuservalue, setGroupuservalue] = useState('');
  const [searchedvalue, setSearchedvalue] = useState('');
  const [isDrawerOpen, setIsDrawerOpen] = useState(true);
  const [searchisOpen, setSearchisOpen] = useState(false);
  const [useraddPopup, setUseraddpopup] = useState(false);
  const [groupUserPopup, setGroupUserPopup] = useState(false);
  const [addeduser, setAddeduser] = useState([]);
  const [groupadduser, setGroupAddeduser] = useState([])
  const [openModal, setOpenModal] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  const [selectedChannel, setSelectedChannel] = useState({});
  const [loadingChat, setLoadingChat] = useState(false);
  const [total, setTotal] = useState(null);
  const [leaveChannelConfirmation, setLeaveChannelConfirmation] = useState({
    open: false,
    data: {},
    type: '',
  });
  const [idle, setIdle] = useState(false);
  const { sendMessage, lastMessage, readyState } = useContext(WebSocketContext);


  useEffect(() => {
    if (readyState === ReadyState.OPEN) {
      setIsConnected(true)
      dispatch(handlers.getDirectChannels());
      dispatch(handlers.getGroupChannels());
      dispatch(handlers.getFavouritesChannels());
      dispatch(handlers.getAllUsers());
    }
  }, [readyState, dispatch]);

  // listing to websocket events
  useEffect(() => {
    if (readyState === ReadyState.OPEN && lastMessage !== null) {
      const message = JSON.parse(lastMessage.data);

      // this event will be triggered when user memebers are added to the channel
      if (message.action === WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_USERS) {
        if (message.target === selectedChannelState?._id) {
          dispatch(handlers.joinGroupChannel(message.target))
        }
        dispatch({
          type: types.UPDATE_CHANNEL_STATE_ON_ADD_USER,
          payload: {
            channel_id: message.target,
            users: message.payload,
          }
        })

        // this event will be triggered when any user leaves the channel
      } else if (message.action === WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_CHANNEL) {
        if (selectedChannelState?._id === message.target) {
          dispatch(handlers.joinGroupChannel(message.target))
        }

        dispatch({
          type: types.UPDATE_USERS_COUNT_ON_LEAVE_CHANNEL,
          payload: {
            channel_id: message.target,
            user_id: message.payload,
          }
        })

        // this event will be triggered when any user is added to group
      } else if (message.action === WEBSOCKET_ACTIONS.ADD_CHANNEL_ON_ADDING_MEMBER) {
        dispatch({
          type: types.ADD_CHANNEL_ON_ADD_IN_GROUP,
          payload: message.payload,
        })
      }
    }
  }, [readyState, dispatch, lastMessage])

  const openSearchUser = () => {
    fetchSearchedChannel()
    setSearchisOpen(true);
  };

  const closeSearchUser = () => {
    setSearchisOpen(false);
  };

  const openAddUserPopupinChannel = (channelId) => {
    setUseraddpopup(true);
    setSelectedChannel(channelId);
  };

  const openGroupUserPopupinChannel = () => {
    setGroupUserPopup(true);
  };

  const closeGroupUserPopupinChannel = () => {
    setGroupUserPopup(false);
    setGroupAddeduser([]);
  }

  const closeAddUserPopupinChannel = () => {
    setUseraddpopup(false);
    setSelectedChannel({})
    setAddeduser([]);
  };

  const handleChangeModal = () => {
    setOpenModal(!openModal);
  };

  const handleClose = () => {
    setLeaveChannelConfirmation({
      open: false,
      data: {},
      type: '',
    });
  };
  
  const handleLeaveChannnel = (data, type) => {
    setLeaveChannelConfirmation({
      open: true,
      data: data,
      type: type,
    });
  };


  // this function will be triggered when user clicks on leave channel
  const leaveChannel = async () => {
    try {

      // this api will be triggered to remove the user from the channel
      const result = await api.put('v1/channel/leavechannel', {
        channel_id: leaveChannelConfirmation?.data?._id,
      });
      if (result?.data?.meta.code === API_STATUS.SUCCESS_CODE) {

        // if the channel is group channel is in the group section then remove the channel from the group section
        if (leaveChannelConfirmation?.type === SECTIONS.GROUP_CHANNEL_SECTION) {
          dispatch({
            type: types.LEAVE_CHANNEL_GROUP,
            payload: leaveChannelConfirmation?.data,
          })

          // if the channel is group channel is in the favourite section then remove the channel from the favourite section
        } else {
          dispatch({
            type: types.LEAVE_CHANNEL_FAVOURITE,
            payload: leaveChannelConfirmation?.data,
          })
        }
      } else {

        // close the leave channel popup and show the error message
        setLeaveChannelConfirmation({
          open: false,
          data: {},
          type: '',
        });

        // show the error message
        toast.error(
          result?.data?.meta.message,
          toastOptions
        );
        return;
      }

      // the userIds to which we have to send the notification
      var userIds = []
      for (let i = 0; i < leaveChannelConfirmation?.data?.users?.length; i++) {
        if (leaveChannelConfirmation.data.users[i].user_id !== userId)
          userIds.push(leaveChannelConfirmation.data.users[i].user_id)
      }

      // this event will be triggered when user leaves the channel
      const WebSocketPayload = {
        action: WEBSOCKET_ACTIONS.POST_GLOBAL_NOTIFICATION_TO_USERS,
        notification_type: WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_CHANNEL,
        target: leaveChannelConfirmation?.data._id,
        payload: userId,
        notify_users: userIds,
      };
      const data = JSON.stringify(WebSocketPayload)
      sendMessage(data)

      // if the user is in the channel then remove the channel from the selected channel
      if (selectedChannelState?._id === leaveChannelConfirmation?.data._id) {
        dispatch({
          type: types.REMOVE_SELECTED_CHANNEL,
        })
      }

      // close the leave channel popup
      setLeaveChannelConfirmation({
        open: false,
        data: {},
        type: '',
      });

    } catch (err) {
      console.log('Error: ', err)
    }
  };


  // this function will be triggered when user clicks on logout
  const handleLogout = () => {
    localStorage.clear();
    window.location.reload();
  };


  // this function will be triggered when user want to search the user based on search value or get all the users
  const fetchSearchedUser = useCallback(async () => {
    try {

      // this api will be triggered to get all the users based on search value or get all the users
      const result = await api.get(`/v1/chat/get-all-users?searchValue=${searchuservalue}`);

      if (result?.data.meta.code === API_STATUS.SUCCESS_CODE) {
        const data = result?.data?.data || [];

        // to filter out the user which are already added in the channel
        const filteredUsers = await data.filter(
          (user) => !addeduser?.some((addedUser) => addedUser._id === user._id)
        );

        // set the filtered users in the state
        dispatch({
          type: types.SET_ALL_USERS,
          payload: filteredUsers,
        })
      } else {
        // if the api fails then set the empty array in the state
        dispatch({
          type: types.SET_ALL_USERS,
          payload: [],
        })
      }
    } catch (err) {
      console.log('Error : ', err);
    }
  }, [searchuservalue, addeduser])


  // this function will be triggered when user want to search the users based on search value or get all the users
  const fetchSearchedGroupUser = useCallback(async () => {
    try {
      const result = await api.get(`/v1/chat/get-all-users?searchValue=${searchgroupuservalue}`);

      // this api will be triggered to get all the users based on search value or get all the users
      if (result?.data.meta.code === API_STATUS.SUCCESS_CODE) {
        const data = result?.data?.data || [];
        // to filter out the user which are already added in the channel
        const filteredUsers = await data.filter(
          (user) => !groupadduser?.some((groupUser) => groupUser._id === user._id)
        );

        // set the filtered users in the state
        dispatch({
          type: types.SET_ALL_USERS,
          payload: filteredUsers,
        })
      } else {
        // if the api fails then set the empty array in the state
        dispatch({
          type: types.SET_ALL_USERS,
          payload: [],
        })
      }
    } catch (err) {
      console.log('Error : ', err);
    }
  }, [searchgroupuservalue, groupadduser])

  // this function will be triggered when user want to fetch the chat history of selected channel
  const fetchChatHistory = async (channel, page = 1) => {

    // set the loading to true till the chat history is fetched
    setLoadingChat(true)
    try {

      // this api will be triggered to get the chat history of the selected channel with given page and offset
      const result = await api.get(`v1/chat/messages-by-channelid/${channel._id}?page=${page}&offset=11`)

      if (result?.data?.meta.code === API_STATUS.SUCCESS_CODE) {
        let data = result?.data?.data?.data || []
        // to reverse the chat history 
        let newChat = data && [...data].reverse();

        // set the chat history in the state
        dispatch({
          type: types.SET_CHAT_HISTORY,
          payload: newChat,
        })

        // set the page and total in the state
        setPage(result?.data?.data?.meta?.page || 1);
        // set the total in the state
        setTotal(result?.data?.data?.meta?.total_page || 0)
        // set the loading to false
        setLoadingChat(false)
      }
    } catch (err) {
      console.log('Error: ', err)
    }
  }

  // this function will be triggered when user want to search the user's channels based on search value
  const fetchSearchedChannel = useCallback(async () => {
    try {

      // this api will be triggered to get all the user's channels based on search value
      const result = await api.get(`/v1/chat/searchhandler?searchValue=${searchedvalue}`);

      if (result?.data.meta.code === API_STATUS.SUCCESS_CODE) {

        const data = result?.data?.data || [];
        // to filter out the users which except the logged in user
        const filteredData = data.filter((item) => item._id !== userId);
        // set the filtered users in the state
        setSearchedChannel(filteredData);
      } else {
        // if the api fails then set the empty array in the state
        setSearchedChannel([]);
      }
    } catch (err) {
      console.log('Error : ', err);
    }
  }, [searchedvalue]);

  // this function will be triggered when user will select any channel from channel list 
  const handleSelectUser = async (user, value) => {

    // clear the prev state of chat history if any
    dispatch({
      type: types.CLEAR_CHAT_HISTORY,
      payload: [],
    })

    // if the user is selected from search and channel type is direct channel then join the direct channel
    if (value === SELECTION_TYPES.FROM_SEARCH && user.channel_type === CHANNEL_TYPES.DIRECT_CHANNEL) {
      const data = dispatch(handlers.joinDirectChannel(user._id))
      await data.then(function (result) {
        user._id = result[0]._id
      })
    }
    // if the user is selected from search and channel type is group channel then join the group channel
    else if (value === SELECTION_TYPES.FROM_SEARCH && user.channel_type === CHANNEL_TYPES.GROUP_CHANNEL) {
      dispatch(handlers.joinGroupChannel(user._id))
    }
    // if the user is selected from channel list then set the selected channel in the state
    else if (value === SELECTION_TYPES.FROM_LIST) {
      dispatch({
        type: types.SELECTED_CHANNEL,
        payload: user,
      });
    }

    // this websocket event will be triggered to register the user in room channel of websocket server
    const webSocketPayload = {
      action: WEBSOCKET_ACTIONS.JOIN_CHANNEL,
      message: user._id,
    };
    let data = JSON.stringify(webSocketPayload)
    sendMessage(data)

    // fetch the chat history of selected channel
    fetchChatHistory(user)
  };

  useEffect(() => {
    if (isConnected) fetchSearchedGroupUser();
  }, [isConnected, searchgroupuservalue, groupUserPopup, fetchSearchedGroupUser]);

  useEffect(() => {
    if (isConnected) fetchSearchedUser();
  }, [isConnected, searchuservalue, useraddPopup, fetchSearchedUser]);

  useEffect(() => {
    fetchSearchedChannel();
  }, [fetchSearchedChannel, searchedvalue]);

  useEffect(() => {
    let idleTimer;

    // Check if the WebSocket connection is open
    if (readyState !== ReadyState.OPEN) return;

    // Function to reset the idle timer
    const resetIdleTimer = () => {
      clearTimeout(idleTimer);
      idleTimer = setTimeout(() => {
        // If user is idle, set status to away and send WebSocket message
        if (idle) {
          setIdle(false);

          // Send WebSocket payload to change user status to away
          const webSocketPayload = {
            action: WEBSOCKET_ACTIONS.CHANGE_USER_STATUS,
            status: USER_STATUS.AWAY,
            payload: {
              user_id: userId,
              status: USER_STATUS.AWAY,
            }
          };
          const data = JSON.stringify(webSocketPayload);
          sendMessage(data);
        }
      }, 10000); // 10 seconds
    };

    // Function to handle user activity
    const handleUserActivity = () => {
      if (!idle) {
        // If user becomes active, set status to online and send WebSocket message
        setIdle(true);
        const webSocketPayload = {
          action: WEBSOCKET_ACTIONS.CHANGE_USER_STATUS,
          status: USER_STATUS.ONLINE,
          payload: {
            user_id: userId,
            status: USER_STATUS.ONLINE,
          }
        };
        const data = JSON.stringify(webSocketPayload);
        sendMessage(data);
      }
      resetIdleTimer(); // Reset the idle timer
    };

    // Attach event listeners for user activity
    window.addEventListener('mousemove', handleUserActivity);
    window.addEventListener('keydown', handleUserActivity);

    // Set up initial timer
    resetIdleTimer();

    // Cleanup: remove event listeners and clear timer
    return () => {
      window.removeEventListener('mousemove', handleUserActivity);
      window.removeEventListener('keydown', handleUserActivity);
      clearTimeout(idleTimer);
    };
  }, [idle, readyState, userId, sendMessage]);



  return (
    <>
      <div className='chatParentContainer'>
        {isConnected ? (
          <>

            {/* Channel List */}
            <UserList
              onSelectUser={handleSelectUser}
              handleLeaveChannel={handleLeaveChannnel}
              isOpen={isDrawerOpen}
              OpenAddUserPopup={openAddUserPopupinChannel}
              searchOpen={openSearchUser}
              onLogout={handleLogout}
              handleChangeModal={handleChangeModal}
            />

            {/* Chat Screen */}
            <ChatScreen
              fetchChat={fetchChatHistory}
              isOpen={openGroupUserPopupinChannel}
              page={page}
              setLoading={setLoadingChat}
              userId={userId}
              loading={loadingChat}
              total={total}
            />

            {/* Group Users Popup */}
            {groupUserPopup && (
              <GroupUserPopup
                searchedvalue={searchgroupuservalue}
                setSearchedvalue={setGroupuservalue}
                onClose={closeGroupUserPopupinChannel}
                groupUserPopup={groupUserPopup}
                setAddeduser={setGroupAddeduser}
                addeduser={groupadduser}
              />
            )}

            {/* Search Channels Popup */}
            <SearchUser
              searchedChannel={searchedChannel}
              setSearchedChannel={setSearchedChannel}
              setSearchedvalue={setSearchedvalue}
              searchedvalue={searchedvalue}
              isOpen={openSearchUser}
              onSelectUser={handleSelectUser}
              onClose={closeSearchUser}
              searchisOpen={searchisOpen}
            />

            {/* Create Channel Popup */}
            {openModal && (
              <ChannelModal
                handleChangeModal={handleChangeModal}
              />
            )}

            {/* Leave Channel Popup */}
            {leaveChannelConfirmation?.open && (
              <ConfirmationModal
                isOpen={leaveChannelConfirmation.open}
                onClose={handleClose}
                title={`Leave private channel ${leaveChannelConfirmation?.data?.channel_name}`}
                message={`Are you sure you want to leave the private channel ${leaveChannelConfirmation?.data?.channel_name}? You must be re-invited in order to re-join this channel in the future.`}
                confirmLabel='Yes, leave channel'
                cancelLabel='Cancel'
                onConfirm={leaveChannel}
              />
            )}

          </>) : (
          <Loading />
        )}
      </div>
      <ToastContainer />
    </>
  );
};

const mapStateToProps = state => ({
  selectedChannelState: state.SelectedChannelReducer || {},
});

export default connect(
  mapStateToProps
)(ChatParent);