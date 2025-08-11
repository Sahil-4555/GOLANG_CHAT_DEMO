export const API_STATUS = {
  SUCCESS_CODE: 1,
  FAILURE_CODE: 0
};

export const CHANNEL_TYPES = {
  DIRECT_CHANNEL: 'one-to-one',
  GROUP_CHANNEL: 'private'
}

export const SELECTION_TYPES = {
  FROM_LIST: 'LIST',
  FROM_SEARCH: 'SEARCH',
  ICON: 'ICON',
  DOWNLOAD: 'DOWNLOAD'
};


export const MEDIA_TYPE = {
  IMAGE: 'image',
  VIDEO: 'video'
};

export const CONTENT_TYPE = {
  TEXT: 1,
  IMAGE: 2,
  VIDEO: 3,
  DOCUMENT: 4,
  SYSTEM: 5,
};

export const CHANNEL_OPTIONS = [
  {
    id: 1,
    title: 'Public Channel',
    description: 'Anyone can join'
  },
  {
    id: 2,
    title: 'Private Channel',
    description: 'Only invited members'
  }
];

export const SECTIONS = {
  DIRECT_CHANNEL_SECTION: "DIRECT_CHANNEL_SECTION",
  GROUP_CHANNEL_SECTION: "GROUP_CHANNEL_SECTION",
  FAVOURITE_CHANNEL_SECTION: "FAVOURITE_CHANNEL_SECTION"
}

export const WEBSOCKET_ACTIONS = {
  POST_GLOBAL_NOTIFICATION_TO_USERS: "post-notification-to-user",
  UPDATE_CHANNEL_DATA_ACROSS_USERS: "update-channel-data-across-user",
  JOIN_CHANNEL: "join-room",
  SEND_MESSAGE: "send-message",
  UPDATE_CHANNEL_DATA_ACROSS_CHANNEL: "update-channel-data-across-channel",
  UPDATE_CHANNEL_ON_MESSAGE: "update-channel-on-message",
  UPDATE_MESSAGE: "update-message",
  DELETE_MESSAGE: "delete-message",
  USER_TYPING: "user-typing",
  USER_STOP_TYPING: "user-stop-typing",
  CHANGE_USER_STATUS: "change-user-status",
  ADD_CHANNEL_ON_ADDING_MEMBER: "add-channel-on-adding-member",
}

export const USER_STATUS = {
  ONLINE: 1,
  OFFLINE: 2,
  AWAY: 3,
  DO_NOT_DISTURB: 4,
}

export const MESSAGES = {
  DELETE_MESSAGE: "<p><em>message deleted</em></p>"
}

export const ROUTES = {
  LOGIN: `${process.env.REACT_APP_API_URL}/v1/api/login`,
  SIGNUP: `${process.env.REACT_APP_API_URL}/v1/api/signup`
}



