import { combineReducers } from "redux";
import ChatHistoryReducer from "./reducers/chathistory";
import DirectChannelReducer from "./reducers/directchannel";
import FavouriteChannelReducer from "./reducers/favouritechannel";
import GroupChannelReducer from "./reducers/groupchannel";
import SelectedChannelReducer from "./reducers/selectedchannel";
import UsersReducer from "./reducers/users";

const rootReducer = combineReducers({
  ChatHistoryReducer,
  DirectChannelReducer,
  FavouriteChannelReducer,
  GroupChannelReducer,
  SelectedChannelReducer,
  UsersReducer,
})

export default rootReducer;