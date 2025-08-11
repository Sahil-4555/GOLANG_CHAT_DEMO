import * as types from "../../utils/constant/ActionTypes";

const initialState = {
    groupChannels: [],
};

const GroupChannelReducer = (state = initialState.groupChannels, action) => {
    switch (action.type) {
        case types.GET_GROUP_CHANNELS:
            return [...action.payload];

        case types.ADD_TO_FAVOURITES_GROUP:
            return state.filter(channel => channel._id !== action.payload._id);

        case types.REMOVE_FROM_FAVOURITES_GROUP:
            return [...state, action.payload].sort((a, b) =>
                new Date(b.last_activity[0]) - new Date(a.last_activity[0])
            );

        case types.HANDLE_CHANNEL_ON_NEW_MESSAGE:
            var len = state.length
            const channel = state.filter((channels) => channels._id === action.payload._id)
            state = state.filter((channels) => channels._id !== action.payload._id)
            if (state.length < len) {
                state = channel.concat(state)
            }
            return state

        case types.LEAVE_CHANNEL_GROUP:
            return state.filter(channel => channel._id !== action.payload._id);

        case types.ADD_NEW_MESSAGE:
            return state.map(channel =>
                channel._id === action.payload.channel_id
                    ? { ...channel, last_activity: [action.payload.payload[0].updated_at] }
                    : channel
            );

        case types.SELECTED_CHANNEL:
            return state.map(channel =>
                channel._id === action.payload._id
                    ? { ...channel, message_count: 0 }
                    : channel
            );

        case types.ADD_COUNTER_ON_NEW_MESSAGE:
            return state.map(channel =>
                channel._id === action.payload
                    ? { ...channel, message_count: channel.message_count + 1 }
                    : channel
            );

        case types.UPDATE_USER_STATUS:
            return state.map(channel => ({
                ...channel,
                users: channel.users.map(user =>
                    user.user_id === action.payload.user_id
                        ? { ...user, status: action.payload.status }
                        : user
                )
            }));

        case types.UPDATE_USERS_COUNT_ON_ADD_IN_GROUP:
            return state.map(channel =>
                channel._id === action.payload.channel_id
                    ? { ...channel, users: action.payload.users }
                    : channel
            );

        case types.UPDATE_USERS_COUNT_ON_LEAVE_CHANNEL:
            return state.map(channel =>
                channel._id === action.payload.channel_id
                    ? { ...channel, users: channel.users.filter(user => user.user_id !== action.payload.user_id) }
                    : channel
            );

        case types.UPDATE_CHANNEL_STATE_ON_ADD_USER:
            return state.map(channel =>
                channel._id === action.payload.channel_id
                    ? { ...channel, users: [...channel.users, action.payload.users] }
                    : channel
            );

        case types.ADD_CHANNEL_ON_ADD_IN_GROUP:
            return [action.payload, ...state]

        default:
            return state;
    }
};

export default GroupChannelReducer;
