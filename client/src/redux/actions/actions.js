import * as types from '../../utils/constant/ActionTypes'
import api from '../../utils/api';
import * as constants from '../../utils/constant/Constant';
import { useSelector } from 'react-redux';


export const getDirectChannels = () => async (dispatch) => {
    try {
        const result = await api.get('v1/channel/get-direct-channel');
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            const data = result?.data.data || [];
            dispatch({
                type: types.GET_DIRECT_CHANNELS,
                payload: data,
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Get Direct Channel Response",
        });
    }
}

export const getGroupChannels = () => async (dispatch) => {
    try {
        const result = await api.get('v1/channel/get-group-channel');
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            const data = result?.data.data || [];
            dispatch({
                type: types.GET_GROUP_CHANNELS,
                payload: data,
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Get Group Channels Error",
        });
    }
}

export const getFavouritesChannels = () => async (dispatch) => {
    try {
        const result = await api.get('v1/channel/get-favourite-channel');
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            const data = result?.data.data || [];
            dispatch({
                type: types.GET_FAVOURITES_CHANNELS,
                payload: data,
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Get Favourites Channels Error",
        });
    }
}

export const handlePrivateAddFavourites = (channel) => async (dispatch) => {
    try {

        dispatch({
            type: types.ADD_TO_FAVOURITES_PRIVATE,
            payload: channel,
        })

        const result = await api.post('v1/channel/addfavoritechannel', {
            is_favourite: true,
            channel_id: channel._id,
        });

        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {

        }
        else {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Private Add Favourite Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Private Add Favourite Error",
        });
    }
}

export const handleGroupAddFavourites = (channel) => async (dispatch) => {
    try {
        dispatch({
            type: types.ADD_TO_FAVOURITES_GROUP,
            payload: channel
        })

        const result = await api.post('v1/channel/addfavoritechannel', {
            is_favourite: true,
            channel_id: channel._id,
        });

        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {

        }
        else {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Private Add Favourite Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Private Add Favourite Error",
        });
    }
}

export const handleRemovePrivateFavourites = (channel) => async (dispatch) => {
    try {
        dispatch({
            type: types.REMOVE_FROM_FAVOURITES_PRIVATE,
            payload: channel,
        })

        const result = await api.post('v1/channel/addfavoritechannel', {
            is_favourite: false,
            channel_id: channel._id,
        });
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {

        }
        else {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Private Add Favourite Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Private Add Favourite Error",
        });
    }
}

export const handleRemoveGroupFavourites = (channel) => async (dispatch) => {
    try {
        dispatch({
            type: types.REMOVE_FROM_FAVOURITES_GROUP,
            payload: channel,
        })

        const result = await api.post('v1/channel/addfavoritechannel', {
            is_favourite: false,
            channel_id: channel._id,
        });
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {

        }
        else {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Private Add Favourite Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Private Add Favourite Error",
        });
    }
}

export const handleCloseConversation = (channel, section, selectedChannelId) => async (dispatch) => {

    try {
        const result = await api.put('v1/channel/closeconversation', {
            channel_id: channel._id,
        });
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            if (section === constants.SECTIONS.DIRECT_CHANNEL_SECTION) {
                dispatch({
                    type: types.CLOSE_CONVERSATION_DIRECT_CHANNEL,
                    payload: channel,
                })
            } else {
                dispatch({
                    type: types.CLOSE_CONVERSATION_FAVOURITE_CHANNEL,
                    payload: channel
                })
            }

            if (selectedChannelId === channel._id) {
                dispatch({
                    type: types.REMOVE_SELECTED_CHANNEL,
                });
            }

        }
        else {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Private Add Favourite Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Private Add Favourite Error",
        });
    }
}

export const handleLeaveChannel = (channel, section) => async (dispatch) => {

    try {
        const result = await api.put('v1/channel/leavechannel', {
            channel_id: channel._id,
        });
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            if (section === constants.SECTIONS.GROUP_CHANNEL_SECTION) {
                dispatch({
                    type: types.LEAVE_CHANNEL_GROUP,
                    payload: channel,
                })
                dispatch(getGroupChannels())
            } else {
                dispatch({
                    type: types.LEAVE_CHANNEL_FAVOURITE,
                    payload: channel,
                })
                dispatch(getFavouritesChannels())
            }
        }
        else {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Private Add Favourite Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Private Add Favourite Error",
        });
    }
}

export const getAllUsers = () => async (dispatch) => {

    try {
        const result = await api.get('v1/chat/get-all-users');
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            const data = result?.data.data || [];
            dispatch({
                type: types.SET_ALL_USERS,
                payload: data,
            })
        }
        else {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Get All Users Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Get All Users Error",
        });
    }
}

export const joinDirectChannel = (userId) => async (dispatch) => {

    try {
        const result = await api.post('v1/channel/joinRoom', {
            reciver_id: userId,
        });
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            dispatch({
                type: types.SELECTED_CHANNEL,
                payload: result.data.data[0],
            })
            // setSelectedUserId(result.data.data._id)
            dispatch(getDirectChannels())
            dispatch(getFavouritesChannels())
            return result?.data?.data;
        }
    } catch (err) {
        console.log('Error: ', err)
    }
}

export const joinGroupChannel = (channelId) => async (dispatch) => {

    try {
        const result = await api.get(`v1/channel/join-group/${channelId}`)
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            dispatch({
                type: types.SELECTED_CHANNEL,
                payload: result.data.data[0],
            })

        }
    } catch (err) {
        console.log('Error: ', err)
    }
}

export const handleNewMessage = (channelId, message) => async (dispatch) => {

    const channelState = useSelector((state) => state);
    try {
        const result = await api.get(`v1/channel/join-group/${channelId}`)
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            dispatch({
                type: types.SELECTED_CHANNEL,
                payload: result.data.data[0],
            })

            if (channelState?.SelectedChannelReducer?._id === channelId) {
                dispatch({
                    action: types.ADD_NEW_MESSAGE,
                    payload: message
                });
            }

        }
    } catch (err) {
        console.log('Error: ', err)
    }
}

export const addChannelOnNewMessage = (channelId) => async (dispatch) => {

    try {
        const result = await api.get(`v1/channel/join-group/${channelId}`)
        if (result?.data?.meta.code === constants.API_STATUS.SUCCESS_CODE) {
            dispatch({
                type: types.ADD_CHANNEL_ON_NEW_MESSAGE,
                payload: result.data.data[0],
            })
        }
    } catch (err) {
        console.log('Error: ', err)
    }
}

export const HandleEditMessage = (messageId, content) => async (dispatch) => {
    try {
        dispatch({
            type: types.EDIT_MESSAGE,
            payload: {
                messageId,
                content,
            },
        })
        const result = await api.put('v1/chat/updatemessage', {
            _id: messageId,
            content: content,
        });
        if (result?.data?.meta.code !== constants.API_STATUS.SUCCESS_CODE) {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Edit Message Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Edit Message Error",
        });
    }
}

export const HandleDeleteMessage = (messageId) => async (dispatch) => {
    try {
        dispatch({
            type: types.DELETE_MESSAGE,
            payload: {
                messageId,
            },
        })
        const result = await api.put('v1/chat/deletemessage', {
            _id: messageId,
        });
        if (result?.data?.meta.code !== constants.API_STATUS.SUCCESS_CODE) {
            dispatch({
                type: types.ERROR_GETTING_CHANNEL,
                payload: "Handle Edit Message Error",
            });
        }
    } catch (error) {
        console.log('Error : ', error);
        dispatch({
            type: types.ERROR_GETTING_CHANNEL,
            payload: "Handle Edit Message Error",
        });
    }
}





