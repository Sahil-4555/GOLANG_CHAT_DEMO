import * as types from "../../utils/constant/ActionTypes"

const initalState = {
    users: [],
}

const UsersReducer = (state = initalState.users, action) => {
    switch (action.type) {
        case types.SET_ALL_USERS:
            return [...action.payload]

        case types.UPDATE_USER_STATUS:
            return state.map(user =>
                user._id === action.payload.user_id
                    ? { ...user, status: action.payload.status }
                    : user
            );

        default:
            return state
    }
}

export const getUsers = (state) => state.UsersReducer;

export default UsersReducer;
