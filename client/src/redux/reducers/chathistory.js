import * as types from "../../utils/constant/ActionTypes";

const initialState = {
    chatHistory: [],
};

const ChatHistoryReducer = (state = initialState.chatHistory, action) => {
    switch (action.type) {
        case types.SET_CHAT_HISTORY:
            return [
                ...action.payload,
                ...state,
            ];

        case types.ADD_NEW_MESSAGE:
            return [
                ...state,
                ...action.payload.payload,
            ];

        case types.CLEAR_CHAT_HISTORY:
            return [...action.payload];

        case types.EDIT_MESSAGE:
            return state.map(message =>
                message._id === action.payload.messageId
                    ? { ...message, content: action.payload.content, is_edited: true }
                    : message
            );

        case types.DELETE_MESSAGE:
            return state.map(message =>
                message._id === action.payload.messageId
                    ? { ...message, is_deleted: true }
                    : message
            );

        case types.UPDATE_USER_STATUS:
            return state.map(message =>
                message.user.user_id === action.payload.user_id
                    ? { ...message, user: { ...message.user, status: action.payload.status } }
                    : message
            );

        default:
            return state;
    }
};

export default ChatHistoryReducer;
