import * as types from "../../utils/constant/ActionTypes"

const initalState = {
    selectedChannel: {},
}

const SelectedChannelReducer = (state = initalState.selectedChannel, action) => {
    switch (action.type) {
      
        case types.SELECTED_CHANNEL:
            return action.payload
        case types.REMOVE_SELECTED_CHANNEL:
            return {}
        default:
            return state
    }
}

export default SelectedChannelReducer;
