/* eslint-disable no-unused-vars */
import React, { Fragment, useCallback, useContext, useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
import 'react-quill/dist/quill.snow.css';
import 'quill-emoji/dist/quill-emoji.css';
import * as Emoji from 'quill-emoji';
import QuillMention from 'quill-mention';
import 'quill-mention/dist/quill.mention.css';
import { ArrowDownIcon } from '@heroicons/react/20/solid';
import * as types from '../utils/constant/ActionTypes'
import NewEditor from '../components/NewEditor';
import WebSocketContext from '../utils/context/WebSocketContext';
import * as constants from '../utils/constant/Constant';
import { connect, useDispatch } from 'react-redux';
import { ReadyState } from 'react-use-websocket';
import * as handlers from '../redux/actions/actions';
import api from '../utils/api';
import axios from 'axios';
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import editMessageIcon from '../assets/icons/editIcon.svg'
import deleteMessageIcon from '../assets/icons/deleteIcon.svg'
import userIcon from '../assets/icons/user.svg'
import submitIcon from '../assets/icons/submitIcon.svg'
import DefaultChatScreen from '../components/DefaultChatScreen';
import Video from '../components/Vedio';
import ConfirmationModal from '../components/ActionPopup';

const ChatScreen = ({
    isOpen,
    fetchChat,
    page,
    total,
    loading,
    setLoading,
    userId,
    directChannelState,
    favouriteChannelState,
    chatHistoryState,
    selectedChannelState
}) => {
    const userName = localStorage.getItem("userName");
    const name = localStorage.getItem("name")
    const messagesEndRef = useRef(null);
    const typingTimeoutRef = useRef(null);
    const editRef = useRef(null);
    const [typedMessage, setTypedMessage] = useState('');
    const [showError, setShowError] = useState({ image: false, file: false });
    const [editMessage, setEditMessage] = useState('');
    const [selectedMessage, setSelectedMessage] = useState({});
    const [fetchedPages, setFetchedPages] = useState([]);
    const [openPopConfirm, setOpenPopConfirm] = useState(false);
    const [selfTyping, setSelfTyping] = useState(false);
    const [typingData, setTypingData] = useState({
        name: '',
        channel: ''
    });
    const [selectedFile, setSelectedFile] = useState({});
    const [imageLoader, setImageLoader] = useState(false);
    const [fileName, setFileName] = useState('');
    const [fileUrl, setFileUrl] = useState('');
    const [uploadProgress, setUploadProgress] = useState(0);
    const [maxSize, setMaxSize] = useState({});
    const [fileType, setFileType] = useState('');
    const [playVideo, setPlayVideo] = useState(false);
    const [videoUrl, setVideoUrl] = useState('');
    const [selectedVideoMessage, setSelectedVideoMessage] = useState({});
    const [mentionPopup, setMentionpopup] = useState(false);
    const inputRef = useRef(null);
    const { sendMessage, lastMessage, readyState } = useContext(WebSocketContext);
    const dispatch = useDispatch();
    const toastOptions = {
        position: "top-right",
        autoClose: 3000,
        pauseOnHover: true,
        draggable: true,
        theme: "dark",
    };

    // Websocket event listener
    useEffect(() => {
        if (readyState === ReadyState.OPEN && lastMessage !== null) {
            const message = JSON.parse(lastMessage.data);

            // this event will be triggered when a new message is sent
            if (message.action === constants.WEBSOCKET_ACTIONS.SEND_MESSAGE) {
                const data = {
                    _id: message._id,
                    content: message.message,
                    file_name: message.file_name,
                    content_type: message.content_type,
                    created_at: message.created_at,
                    updated_at: message.updated_at,
                    is_deleted: false,
                    is_edited: false,
                    user: {
                        name: message.sender.name,
                        status: message.status,
                        user_id: message.sender._id,
                        user_name: message.sender.user_name,
                    },
                }
                const payload = []
                payload.push(data);

                // dispatch the new message to the store to update the chat history
                dispatch({
                    type: types.ADD_NEW_MESSAGE,
                    payload: {
                        payload: payload,
                        channel_id: message.target,
                    },
                })
            }
            else if (message.action === constants.WEBSOCKET_ACTIONS.UPDATE_CHANNEL_ON_MESSAGE) {
                // this event will be triggered when a new message gets to update the channel state in the store

                // if the message is of the channel type direct channel then check if the channel is already there in the direct channel state
                if (message.payload.channel_type === constants.CHANNEL_TYPES.DIRECT_CHANNEL) {
                    let is_channel_there = false
                    for (var i = 0; i < directChannelState?.length; i++) {

                        if (directChannelState[i]._id === message.payload._id) {
                            is_channel_there = true;
                            break
                        }
                    }

                    for (i = 0; i < favouriteChannelState?.length && !is_channel_there; i++) {

                        if (favouriteChannelState[i]._id === message.payload._id) {
                            is_channel_there = true;
                            break;
                        }
                    }

                    // if the channel is not there in the direct channel state & favourite channel state then add the channel to the direct channel state
                    if (!is_channel_there) {
                        dispatch(handlers.addChannelOnNewMessage(message.payload._id))
                    } else {
                        // if the channel is there in the direct channel state then update the channel state
                        dispatch({
                            type: types.HANDLE_CHANNEL_ON_NEW_MESSAGE,
                            payload: message.payload,
                        })
                    }
                } else {
                    // if the channel type is group channel then update the group channel state
                    dispatch({
                        type: types.HANDLE_CHANNEL_ON_NEW_MESSAGE,
                        payload: message.payload,
                    })
                }
                // if currently the selected channel is not the channel on which the message is sent then add the counter to the channel of new message received
                if (selectedChannelState?._id !== message.payload._id) {
                    dispatch({
                        type: types.ADD_COUNTER_ON_NEW_MESSAGE,
                        payload: message.payload._id,
                    })
                }
            }
            else if (message.notification_type === constants.WEBSOCKET_ACTIONS.UPDATE_MESSAGE) {
                // this event will be triggered when any message gets updated in the channel

                // if the selected channel is the channel on which the message is updated then update the message in the chat history
                if (selectedChannelState?._id === message.payload.channel_id) {
                    dispatch({
                        type: types.EDIT_MESSAGE,
                        payload: {
                            messageId: message.payload._id,
                            content: message.payload.content,
                        },
                    })
                }
            }
            else if (message.notification_type === constants.WEBSOCKET_ACTIONS.DELETE_MESSAGE) {
                // this event will be triggered when any message gets deleted in the channel

                // if the selected channel is the channel on which the message is deleted then delete the message from the chat history
                if (selectedChannelState?._id === message.payload.channel_id) {
                    dispatch({
                        type: types.DELETE_MESSAGE,
                        payload: {
                            messageId: message.payload._id,
                        },
                    })
                }
            }
            else if (message.notification_type === constants.WEBSOCKET_ACTIONS.USER_TYPING) {
                // this event will be triggered when any user starts typing in the channel

                // if the selected channel is the channel on which the user is typing then show the typing indicator
                if (selectedChannelState?._id === message.payload.channel_id) {
                    let data = {
                        user_name: message.payload.name,
                        channel_id: message.payload.channel_id,
                    }
                    handleOnWebsocketTyping(data)
                }
            }
            else if (message.notification_type === constants.WEBSOCKET_ACTIONS.USER_STOP_TYPING) {
                // this event will be triggered when any user stops typing in the channel

                // if the selected channel is the channel on which the user stops typing then hide the typing indicator
                if (selectedChannelState?._id === message?.payload?.channel_id) {
                    setTypingData({})
                }
            }
            else if (message.action === constants.WEBSOCKET_ACTIONS.CHANGE_USER_STATUS) {
                // this event will be triggered when any user changes the status

                // update the user status in the store of the user
                dispatch({
                    type: types.UPDATE_USER_STATUS,
                    payload: message.payload,

                });
            }
        }
    }, [readyState, dispatch, lastMessage])

    // this function will be triggered when any user starts typing in the channel to the set the user typing indicator
    const handleOnWebsocketTyping = (data) => {
        setTypingData({
            name: data?.user_name || '',
            channel: data?.channel_id || ''
        });
    };

    // Debounce function to limit the number of times an expensive function is called
    const debounce = (func, delay) => {
        let timeoutId;
        return function (...args) {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => func.apply(this, args), delay);
        };
    };

    // this function will be triggered when the user sends the message
    const debouncedSendMessage = useCallback(
        debounce((message, file, selectedFile) => {
            const type = selectedFile?.type && selectedFile.type.split('/')[0];
            const content_type =
                type === constants.MEDIA_TYPE.IMAGE
                    ? constants.CONTENT_TYPE.IMAGE
                    : type === constants.MEDIA_TYPE.VIDEO
                        ? constants.CONTENT_TYPE.VIDEO
                        : constants.CONTENT_TYPE.DOCUMENT;

            // websocket payload to send the message
            let websocket_payload = {
                action: constants.WEBSOCKET_ACTIONS.SEND_MESSAGE,
                message: message.trim(),
                status: 1,
                target: selectedChannelState?._id,
                content_type: file && file.length ? content_type : constants.CONTENT_TYPE.TEXT,
                file_name: file
            };
            let data = JSON.stringify(websocket_payload)
            if (readyState === ReadyState.OPEN) sendMessage(data);


            // collect the userid of users that are in the selected channel
            var userIds = []
            for (let i = 0; i < selectedChannelState?.users.length; i++) {
                userIds.push(selectedChannelState?.users[i]?.user_id)
            }

            // websocket payload to send the notification to the users of the selected channel
            let payload = selectedChannelState;
            if (selectedChannelState?.channel_type === constants.CHANNEL_TYPES.DIRECT_CHANNEL) {
                userIds.push(userId)
                payload = {
                    _id: selectedChannelState?._id,
                    channel_type: selectedChannelState?.channel_type,
                    channel_name: name,
                    users: [
                        {
                            user_id: userId,
                            name: name,
                            status: 1,
                            user_name: userName,
                        }
                    ],
                }
            }

            // websocket payload to send the update the channel states in the store
            websocket_payload = {
                action: constants.WEBSOCKET_ACTIONS.POST_GLOBAL_NOTIFICATION_TO_USERS,
                notification_type: constants.WEBSOCKET_ACTIONS.UPDATE_CHANNEL_ON_MESSAGE,
                notify_users: userIds,
                payload: payload,
            }
            data = JSON.stringify(websocket_payload)
            if (readyState === ReadyState.OPEN) sendMessage(data);

            // set the typed message to empty after sending the message
            setTypedMessage('');
            // set the selected file to empty after sending the message
            setSelectedFile({});
            // set teh file url to empty after sending the message
            setFileName('');
            // set the file url to empty after sending the message
            setFileUrl('');
        }, 200),
        [selectedChannelState?._id]
    );

    // create a temporary div element
    function cleanHTML(html) {

        // set the innerHTML of the div element to the html content
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = html;
        tempDiv.innerHTML = tempDiv.innerHTML.replace(/\s{2,}/g, ' ');
        tempDiv.querySelectorAll(':empty').forEach((element) => {
            element.parentNode.removeChild(element);
        });
        tempDiv.querySelectorAll('*').forEach((element) => {
            if (element.innerHTML.trim() === '') {
                element.parentNode.removeChild(element);
            }
        });
        return tempDiv.innerHTML;
    }

    // this function will be triggered when the user sends the message
    const SendMessage = (e) => {
        // prevent the default behaviour of the form
        e.preventDefault();

        // set the self typing to false and send the websocket payload to stop the typing indicator
        let websocket_payload = {
            action: constants.WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_CHANNEL,
            notification_type: constants.WEBSOCKET_ACTIONS.USER_STOP_TYPING,
            target: selectedChannelState?._id,
            payload: {
                name: userName,
                channel_id: selectedChannelState?._id,
            },
        }
        let data = JSON.stringify(websocket_payload)
        if (readyState === ReadyState.OPEN) sendMessage(data);

        // clean the html content of the message
        const message = cleanHTML(typedMessage);
        // send the message to the debouncedSendMessage function
        debouncedSendMessage(message, fileName, selectedFile);
    };

    // this function will be triggered when the user types the message to handle the change in the message
    const handleChangeMessage = (event) => {

        // set the typed message to the event value
        setTypedMessage(event);

        // if the user is typing then send the websocket payload to show the typing indicator
        if (!selfTyping && readyState === ReadyState.OPEN) {
            // set the self typing to true
            setSelfTyping(true);

            // send the websocket payload to show the typing indicator
            let websocket_payload = {
                action: constants.WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_CHANNEL,
                notification_type: constants.WEBSOCKET_ACTIONS.USER_TYPING,
                target: selectedChannelState?._id,
                payload: {
                    name: userName,
                    channel_id: selectedChannelState?._id,
                },
            }
            let data = JSON.stringify(websocket_payload)
            if (readyState === ReadyState.OPEN) sendMessage(data);
        }

        // if the user is typing then clear the timeout and set the timeout to show the typing indicator
        if (typingTimeoutRef?.current) {
            clearTimeout(typingTimeoutRef.current);
        }

        // set the timeout to show the typing indicator
        typingTimeoutRef.current = setTimeout(() => {
            // send the websocket payload to stop the typing indicator
            let websocket_payload = {
                action: constants.WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_CHANNEL,
                notification_type: constants.WEBSOCKET_ACTIONS.USER_STOP_TYPING,
                target: selectedChannelState?._id,
                payload: {
                    name: userName,
                    channel_id: selectedChannelState?._id,
                },
            }
            let data = JSON.stringify(websocket_payload)
            if (readyState === ReadyState.OPEN) sendMessage(data);
            // set the self typing to false
            setSelfTyping(false);
        }, 3000);
    };

    // this function will be triggered when the user edits the message to handle the change in the message
    const handleChangeEditMessage = (event) => {
        setEditMessage(event);
    };

    const handleClose = () => {
        setOpenPopConfirm(false);
        setSelectedMessage({});
      };

    // this function will be triggered when the user submits the edited message
    const editSubmit = (e) => {
        if (e.keyCode === 13 && !e.shiftKey && cleanHTML(editMessage)) {
            if (mentionPopup) {
                setMentionpopup(false);
                return;
            }
            e.preventDefault();

            setEditMessage('');
            setSelectedMessage({});
        }
    };

    // this function will be triggered when the user scrolls the chat screen to the bottom 
    const scrollToBottom = () => {
        if (messagesEndRef && messagesEndRef.current) {
            messagesEndRef?.current?.scrollIntoView();
        }
    };

    // format the timestamp to the time
    const formatTime = (timestamp) => {

        // get the message date
        const messageDate = new Date(timestamp);
        const currentDate = new Date();

        // check if the message date is today
        const isToday =
            messageDate.getDate() === currentDate.getDate() &&
            messageDate.getMonth() === currentDate.getMonth() &&
            messageDate.getFullYear() === currentDate.getFullYear();

        // format the time
        const options = {
            hour: 'numeric',
            minute: 'numeric',
            hour12: true
        };

        // create a new date time format
        const formatter = new Intl.DateTimeFormat('en-US', options);

        // if the message date is today then return the time in the format of hh:mm AM/PM
        if (isToday) {
            return (
                formatter.formatToParts(messageDate).find((part) => part.type === 'hour').value +
                ':' +
                formatter.formatToParts(messageDate).find((part) => part.type === 'minute').value +
                ' ' +
                formatter.formatToParts(messageDate).find((part) => part.type === 'dayPeriod').value
            );
        } else {
            // if the message date is not today then return the date in the format of dd/mm/yyyy
            return formatter.format(messageDate);
        }
    };

    useEffect(() => {
        scrollToBottom();
    }, [chatHistoryState]);

    // this function will be triggered to the message date display in the chat screen
    const getMessageDateDisplay = (createdAt) => {
        const messageDate = new Date(createdAt);
        const currentDate = new Date();
        const yesterday = new Date();
        yesterday.setDate(currentDate.getDate() - 1);
        const differenceInDays = (currentDate - messageDate) / (1000 * 60 * 60 * 24);

        // check if the message date is today
        if (messageDate.toLocaleDateString() === currentDate.toLocaleDateString()) {
            return 'Today';
            // check if the message date is yesterday
        } else if (messageDate.toLocaleDateString() === yesterday.toLocaleDateString()) {
            return 'Yesterday';
            // check if the message date is less than 7 days
        } else if (differenceInDays < 7) {
            const options = { weekday: 'long' };
            return messageDate.toLocaleDateString(undefined, options);
        } else {
            return messageDate.toLocaleDateString();
        }
    };

    // this function is to handle the edit message submit
    const handleSubmitEditMessage = (e) => {
        e.preventDefault();
        // dispatch the edit message to the store to update the chat history    
        dispatch(handlers.HandleEditMessage(selectedMessage._id, editMessage))

        // websocket paylod to update the message to edited message
        let websocket_payload = {
            action: constants.WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_CHANNEL,
            notification_type: constants.WEBSOCKET_ACTIONS.UPDATE_MESSAGE,
            target: selectedChannelState?._id,
            payload: {
                _id: selectedMessage._id,
                channel_id: selectedChannelState?._id,
                content: editMessage,
            },
        }
        let data = JSON.stringify(websocket_payload)
        if (readyState === ReadyState.OPEN) sendMessage(data);

        setEditMessage('');
        setSelectedMessage({});
    };

    const handleOpenPopConfirm = (message) => {
        setOpenPopConfirm(true);
        setSelectedMessage(message);
    };

    // this function is handle to the delete the message
    const handleDeleteMessage = () => {
        // dispatch the delete message to the store to delete the message
        dispatch(handlers.HandleDeleteMessage(selectedMessage._id))

        // websocket paylod to delete the message from the chat history
        let websocket_payload = {
            action: constants.WEBSOCKET_ACTIONS.UPDATE_CHANNEL_DATA_ACROSS_CHANNEL,
            notification_type: constants.WEBSOCKET_ACTIONS.DELETE_MESSAGE,
            target: selectedChannelState?._id,
            payload: {
                _id: selectedMessage._id,
                channel_id: selectedChannelState?._id,
            },
        }
        let data = JSON.stringify(websocket_payload)
        if (readyState === ReadyState.OPEN) sendMessage(data);

        setEditMessage('');
        setSelectedMessage({});
        setOpenPopConfirm(false);
    };

    const handleEditMessage = (message) => {
        const messageValue = message && message?.content;
        setEditMessage(messageValue);
        setSelectedMessage(message);
    };

    // this function is to handle the keyup event
    const handleonKeyup = (e) => {
        if (e.keyCode === 38 && !cleanHTML(typedMessage)) {
            const lastChatBySelf = chatHistoryState
                ?.filter((message) => message.user.user_id === userId)
                .pop();
            if (!lastChatBySelf) return;
            setEditMessage(lastChatBySelf.content);
            setSelectedMessage(lastChatBySelf);
            queueMicrotask(() => {
                editRef?.current?.focus();
            }, 10);
        } else if (e.keyCode === 13 && !e.shiftKey && (cleanHTML(typedMessage) || fileName.length)) {
            if (mentionPopup) {
                setMentionpopup(false);
                return;
            }
            SendMessage(e);
        }
    };

    const handleCancel = () => {
        setEditMessage('');
        setSelectedMessage({});
    };

    useEffect(() => {
        // if the image size is greater than 10MB then show the error message
        if (maxSize.image) {
            setShowError({ image: true, file: false });
            setTimeout(() => setShowError({ image: false, file: false }), 3000);
            // if the file size is greater than 50MB then show the error message
        } else if (maxSize.file) {
            setShowError({ image: false, file: true });
            setTimeout(() => setShowError({ image: false, file: false }), 3000);
        }
    }, [maxSize]);

    useEffect(() => {
        scrollToBottom();
        setOpenPopConfirm(false);
        setSelectedMessage({});
    }, [selectedChannelState?._id]);

    useEffect(() => {
        setTypedMessage('');
        setEditMessage('');
        setFetchedPages([]);
        setSelectedFile({});
        setFileUrl('');
        setFileType('');
        setFileName('');
        setPlayVideo(false);
        setVideoUrl('');
        setSelectedVideoMessage({});
    }, [selectedChannelState?._id]);

    const debouncedFetchChat = debounce((selectedChannel, page) => {
        if (page <= total && !fetchedPages.includes(page)) {
            fetchChat(selectedChannel, page);
            setFetchedPages((prevPages) => [...prevPages, page]);
            setLoading(true);
        }
    }, 400);

    const handleOnScrollCapture = (e) => {
        const scrollableHeight = e.target.scrollHeight;
        const scrollPosition = e.target.scrollTop;
        const isAtBottom = scrollPosition + e.target.clientHeight >= scrollableHeight;

        if (!isAtBottom && scrollPosition < 100) {
            debouncedFetchChat(selectedChannelState, page + 1);
        }
        const scrollBottomThreshold = 50;

        if (
            isAtBottom &&
            scrollableHeight > scrollPosition + e.target.clientHeight + scrollBottomThreshold
        ) {
            scrollToBottom();
            setLoading(false);
        }
    };

    // this function will be triggered when the user clicks on the file upload button
    const handleFileChange = (e) => {
        const input = document.createElement('input');
        input.type = 'file';

        input.click();
        input.onchange = async () => {
            const file = input?.files && input.files[0];

            if (!file) return;

            setImageLoader(true);
            setFileName('');
            setMaxSize({});
            setSelectedFile({});

            const fileUrl = file && URL.createObjectURL(file);
            const typeValue = file?.type && file.type.split('/');
            const type = typeValue && typeValue[0];

            setFileType(type);

            // if the image size is greater than 10MB then show the error message
            if (file?.size && type === constants.MEDIA_TYPE.IMAGE && !(file.size <= 10 * 1024 * 1024)) {
                setMaxSize({ image: true });
                return;
            }

            // if the file size is greater than 50MB then show the error message
            if (file?.size && !(file.size <= 50 * 1024 * 1024)) {
                setMaxSize({ file: true });
                return;
            }

            setSelectedFile(file);

            if (file?.type) {
                const payload = {
                    file_name: file.name,
                    file_size: file.size,
                    channel_id: selectedChannelState?._id,
                };

                setFileUrl(fileUrl);

                // upload the file to the server
                try {
                    // get the pre-signed url to upload the file
                    const response = await api.post('/v1/media/uploadmedia', payload);
                    if (response.data?.meta.code === constants.API_STATUS.FAILURE_CODE) {
                        toast.error(response?.data?.meta?.message, toastOptions);
                        return;
                    }

                    const uploadFilePath = response.data?.data?.pre_signed_url;

                    // upload the file to sever with the pre-signed url
                    await axios.put(uploadFilePath, file, {
                        headers: {
                            'Content-Type': typeValue[1] === 'quicktime' ? `${typeValue[0]}/mp4` : file.type
                        },
                        onUploadProgress: (progressEvent) => {
                            const percentCompleted = Math.round(
                                (progressEvent.loaded * 100) / progressEvent.total
                            );
                            setUploadProgress(percentCompleted);
                            if (percentCompleted === 100) setImageLoader(false);
                        }
                    });

                    // set the file name to the file name
                    setFileName(response?.data?.data?.file_name);

                    if (inputRef && inputRef?.current) inputRef.current.focus();
                } catch (error) {
                    toast.error(`Error: ${error.message}`, toastOptions);
                }
            }
        }
    };

    // this function is to handle the format of the file size
    const formatFileSize = (sizeInBytes) => {
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = parseInt(Math.floor(Math.log(sizeInBytes) / Math.log(1024)));

        if (i === 0) {
            return sizeInBytes + ' ' + sizes[i];
        }

        return (sizeInBytes / Math.pow(1024, i)).toFixed(2) + ' ' + sizes[i];
    };

    // this function is to remove the selected file
    const removeSelectedFile = () => {
        setSelectedFile({});
        setFileUrl('');
        setFileType('');
        setFileName('');
    };

    // this fubnciton is to download the file
    const DownloadFile = async (file, type) => {
        if (file && file?.file_name) {
            try {
                let name = file.file_name.substring(file.file_name.lastIndexOf('/') + 1) || '';

                // replace the attachment url
                let newUrl = replaceAttachmentURL(file?.file_name);

                // if the file name ends with .plain then replace it with .txt
                if (name.endsWith('.plain')) name = name.slice(0, -6) + '.txt';
                // if the file name ends with .quicktime then replace it with .mp4
                if (name.endsWith('.quicktime')) name = name.slice(0, -9) + '.mp4';

                // fetch the file from the server
                const response = await fetch(newUrl);

                // if the response is not ok then throw an error
                if (!response.ok) {
                    throw new Error(`Failed to fetch: ${response.status} ${response.statusText}`);
                }

                // create a blob from the response
                const blob = await response.blob();
                // create a object url from the blob
                const objectUrl = URL.createObjectURL(blob);

                // if the file type is vedio then open the vedio in the vedio player
                if (type === constants.SELECTION_TYPES.ICON && file?.content_type === constants.CONTENT_TYPE.VIDEO) {
                    setVideoUrl(objectUrl);
                    setPlayVideo(true);
                    setSelectedVideoMessage(file);
                    return;
                }

                // create a link element to download the file
                const link = document.createElement('a');
                link.href = objectUrl;
                link.download = name;
                link.target = '_blank';

                // append the link to the body and click the link
                document.body.appendChild(link);
                // click the link
                link.click();
                document.body.removeChild(link);
            } catch (error) {
                console.error('Error fetching or creating Blob:', error.message);
            }
        } else {
            console.error('Invalid file or attachment URL.');
        }
    };

    const replaceAttachmentURL = (url) => {
        const newUrl = url.replace(
            'https://fitness-whitelable.s3.ap-south-1.amazonaws.com',
            'https://d3iwp5wdwtnoby.cloudfront.net'
        );
        return newUrl;
    };

    const renderHTMLContent = (htmlContent) => {
        const withoutLastBr = htmlContent.replace(/<p>(\s|&nbsp;)*<br><\/p>$/, '');
        const sanitizedContent = withoutLastBr.replace(/(\s|&nbsp;)+(?=<\/p>)/g, '');
        return { __html: sanitizedContent };
    };

    return (
        <div className='chatMessages'>

            {/* chat screen header */}
            {selectedChannelState?.users?.length ? (
                <div>
                    {selectedChannelState?.users?.length && (
                        <header className='headerContainer'>
                            {selectedChannelState?.channel_name && (
                                <h2 className='channelName'>
                                    {selectedChannelState?.channel_name}
                                </h2>
                            )}
                            {selectedChannelState?.channel_type !== constants.CHANNEL_TYPES.DIRECT_CHANNEL && (
                                <span
                                    className='channelUsers'
                                    onClick={isOpen}>
                                    <img src={userIcon} alt="user icon" style={{ width: 25, height: 25 }} />
                                    {selectedChannelState?.users?.length}
                                </span>

                            )}
                        </header>)}

                    {/* chat screen */}
                    <div className='chatScreenContainer'>
                        {/* loader */}
                        {loading && (
                            <div className='spinner'></div>
                        )}

                        {/* channel message container */}
                        <div
                            className='chatHistoryContainer'
                            onScrollCapture={(e) => handleOnScrollCapture(e)}>
                            {chatHistoryState && chatHistoryState?.length > 0 ? (
                                chatHistoryState?.map((message, index) => {
                                    const isDifferentDate =
                                        index === 0 ||
                                        new Date(message?.created_at).toLocaleDateString() !==
                                        new Date(chatHistoryState[index - 1].created_at).toLocaleDateString();

                                    const isSameUserAsPrevious =
                                        index > 0 &&
                                        message?.user.user_id === chatHistoryState[index - 1]?.user?.user_id &&
                                        message.content_type !== constants.CONTENT_TYPE.SYSTEM &&
                                        !isDifferentDate

                                    return (
                                        <Fragment key={message?._id}>

                                            {/* Date Display Container */}
                                            <div
                                                className={`${isDifferentDate ? 'relative flex-col' : 'hover:bg-blue-100'
                                                    } dateContainer group`}>
                                                {isDifferentDate && (
                                                    <div className='dateDivider'>
                                                        <span className='dateDividerLine'></span>
                                                        <span className='w-fit mx-2'>
                                                            {getMessageDateDisplay(message?.created_at)}
                                                        </span>
                                                        <span className='dateSpacing'></span>
                                                    </div>
                                                )}

                                                {/* Message Content Container */}
                                                <div
                                                    className={`messageContent ${isDifferentDate && 'hover:bg-blue-100 w-full'
                                                        } `}>
                                                    {!isSameUserAsPrevious && (
                                                        <div>
                                                            <div
                                                                className={`avatarContainer`}>
                                                                {message?.content_type === constants.CONTENT_TYPE.SYSTEM ? (
                                                                    <img
                                                                        alt='mindinventory'
                                                                        className='w-[25px] object-contain'
                                                                    />
                                                                ) : (
                                                                    message?.user?.name.charAt(0)
                                                                )}

                                                                <span
                                                                    className={`userStatus`}>
                                                                    {message?.user?.status === constants.USER_STATUS.ONLINE && (
                                                                        <div className='userStatusOnline'></div>
                                                                    )}
                                                                    {message?.user?.status === constants.USER_STATUS.AWAY && (
                                                                        <div className='userStatusAway'></div>
                                                                    )}
                                                                    {message?.user?.status === constants.USER_STATUS.DO_NOT_DISTURB && (
                                                                        <div className='userStatusDonotDisturb'></div>
                                                                    )}
                                                                    {message?.user?.status === constants.USER_STATUS.OFFLINE && (
                                                                        <div className='userStatusOffline'></div>
                                                                    )}
                                                                </span>

                                                            </div>
                                                        </div>
                                                    )}

                                                    <div className='messageDetails'>
                                                        {!isSameUserAsPrevious && (
                                                            <div className='messageUser'>
                                                                <span className='messageUserName'>
                                                                    {message?.content_type !== constants.CONTENT_TYPE.SYSTEM
                                                                        ? message?.user?.name
                                                                        : 'System'}
                                                                </span>
                                                                <span className='messageTime'>
                                                                    {formatTime(message?.created_at)}
                                                                </span>
                                                            </div>
                                                        )}
                                                        <div className={`messageText`}>
                                                            {isSameUserAsPrevious && message?.content_type !== constants.CONTENT_TYPE.SYSTEM && (
                                                                <span className='messageTimestamp invisible group-hover:visible'>
                                                                    {formatTime(message?.created_at)}
                                                                </span>
                                                            )}

                                                            {/* If the content is text */}
                                                            {message?.content_type === constants.CONTENT_TYPE.TEXT &&
                                                                <div>
                                                                    {message?.content_type !== constants.CONTENT_TYPE.SYSTEM ? (
                                                                        <span className='messageContentText'>
                                                                            {message?.is_deleted === false ? (
                                                                                <>
                                                                                    <div
                                                                                        className='chatMessages'
                                                                                        dangerouslySetInnerHTML={renderHTMLContent(message?.content)}
                                                                                    />
                                                                                    <span
                                                                                        className={`editedMessage ${!message?.is_edited ? 'hidden' : ''
                                                                                            } `}>
                                                                                        {message?.is_edited && 'Edited'}
                                                                                    </span>
                                                                                </>
                                                                            ) : (
                                                                                <div
                                                                                    className='chatMessages'
                                                                                    dangerouslySetInnerHTML={renderHTMLContent(constants.MESSAGES.DELETE_MESSAGE)}
                                                                                />
                                                                            )}
                                                                        </span>
                                                                    ) : (
                                                                        <div className='systemMessage' key={message?._id}>
                                                                            <p className='min-w-max text-gray-600 chatMessages '>
                                                                                {message?.content}
                                                                            </p>
                                                                        </div>
                                                                    )}
                                                                </div>
                                                            }

                                                            {/* If the content is media */}
                                                            {message?.file_name && (
                                                                <div className='mediaContainer'>
                                                                    {message?.is_deleted === true ? (
                                                                        <div
                                                                            className='chatMessages'
                                                                            dangerouslySetInnerHTML={renderHTMLContent(constants.MESSAGES.DELETE_MESSAGE)}
                                                                        />) : (
                                                                        <div>
                                                                            {/* {message?.content_type === constants.CONTENT_TYPE.IMAGE ? (
                                                                            <div className='flex h-[18rem] relative'>
                                                                                <img
                                                                                    src={replaceAttachmentURL(message?.file_name)}
                                                                                    alt='Uploaded File'
                                                                                    className='w-full max-w-[30rem]'
                                                                                />
                                                                                <button
                                                                                    className='absolute top-1.5 right-1.5 cursor-pointer bg-gray-600 text-center  p-1  invisible group-hover:visible rounded rounded-sm'
                                                                                    onClick={() => DownloadFile(message)}>
                                                                                    <ArrowDownIcon className='h-5 w-5 text-gray-400' />
                                                                                </button>
                                                                            </div>
                                                                        ) : ( */}
                                                                            <div className='uploadedFile'>
                                                                                <button>
                                                                                    <svg
                                                                                        xmlns='http://www.w3.org/2000/svg'
                                                                                        fill='none'
                                                                                        viewBox='0 0 24 24'
                                                                                        strokeWidth={1.5}
                                                                                        stroke='blue'
                                                                                        className='w-12'
                                                                                        onClick={() => DownloadFile(message, constants.SELECTION_TYPES.ICON)}>
                                                                                        <path
                                                                                            strokeLinecap='round'
                                                                                            strokeLinejoin='round'
                                                                                            d='M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z'
                                                                                        />
                                                                                        {chatHistoryState?.content_type === constants.CONTENT_TYPE.VIDEO && (
                                                                                            <path
                                                                                                strokeLinecap='round'
                                                                                                strokeLinejoin='round'
                                                                                                transform='translate(6.5, 8.25) scale(0.5)'
                                                                                                d='M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971L6.917 18.307a1.125 1.125 0 01-1.667-.985V5.653z'
                                                                                            />
                                                                                        )}
                                                                                    </svg>
                                                                                </button>
                                                                                <div className='file-details'>
                                                                                    <div className='file-name'>
                                                                                        {message?.file_name.substring(
                                                                                            message?.file_name.lastIndexOf('/') + 1
                                                                                        ) || ''}
                                                                                    </div>
                                                                                    <button
                                                                                        className='downloadButton'
                                                                                        onClick={() => DownloadFile(message)}>
                                                                                        <ArrowDownIcon className='arrowDownIcon' />
                                                                                    </button>
                                                                                </div>
                                                                            </div>
                                                                            {/* )} */}
                                                                        </div>
                                                                    )}
                                                                </div>
                                                            )}
                                                        </div>
                                                    </div>
                                                    {userId === message?.user?.user_id && !message?.is_deleted &&
                                                        message.content_type !== constants.CONTENT_TYPE.SYSTEM && (
                                                            <div className='group'>
                                                                <div className='messageMethods'>
                                                                    <div className='messageOptions'>
                                                                        {message.content_type === constants.CONTENT_TYPE.TEXT && (
                                                                            <div>
                                                                                {/* Edit Message  */}
                                                                                < span onClick={() => handleEditMessage(message)}>
                                                                                    <img src={editMessageIcon} alt="edit message icon" style={{ width: 25, height: 25 }} />
                                                                                </span>
                                                                            </div>
                                                                        )}
                                                                        {/* Delete Message  */}
                                                                        <span onClick={() => handleOpenPopConfirm(message)}>
                                                                            <img src={deleteMessageIcon} alt="delete message icon" style={{ width: 25, height: 25 }} />
                                                                        </span>
                                                                    </div>
                                                                </div>
                                                            </div>
                                                        )}
                                                </div>

                                                {/* Edit Message Popup */}
                                                {selectedMessage?._id === message?._id && !openPopConfirm && (
                                                    <div className='popupEditor'>
                                                        <form onSubmit={handleSubmitEditMessage}>
                                                            <div className='flex'>
                                                                <div className='group'>
                                                                    <span className='editorContainer'></span>
                                                                </div>
                                                                <div className='edit-message-editor editMessage'>
                                                                    <NewEditor
                                                                        handleFileChange={handleFileChange}
                                                                        value={editMessage}
                                                                        handleChangeMessage={handleChangeEditMessage}
                                                                        userName={selectedChannelState?.channel_name}
                                                                        setMentionpopup={setMentionpopup}
                                                                        mentionPopup={mentionPopup}
                                                                        handleonKeyup={editSubmit}
                                                                    />
                                                                    <div className='editorButton'>
                                                                        <button
                                                                            type='submit'
                                                                            disabled={
                                                                                (!editMessage.trim().length > 0 || !cleanHTML(editMessage)) &&
                                                                                !fileName.length
                                                                            }
                                                                            className="saveEditorButton" disabledButton>
                                                                            Save
                                                                        </button>
                                                                        <button
                                                                            type='button'
                                                                            className='cancelEditorButton'
                                                                            onClick={() => handleCancel()}>
                                                                            Cancel
                                                                        </button>
                                                                    </div>
                                                                </div>
                                                            </div>
                                                        </form>
                                                    </div>)}
                                            </div>
                                        </Fragment>
                                    );
                                })
                            ) : (selectedChannelState?.users?.length) ? (
                                <div className='startConversation'>Please start a conversation</div>
                            ) : (<div></div>)}
                            <div ref={messagesEndRef} />
                        </div>
                        {/* To show the name of the user which is typing */}
                        {typingData?.name !== userName
                            && typingData?.channel === selectedChannelState?._id && (
                                <div className='text-md text-gray-700'>{`${typingData?.name} is typing...`}</div>
                            )}

                        {selectedChannelState?._id && (
                            <form onSubmit={SendMessage}>
                                {fileUrl && (
                                    <div className='sendMediaMessage'>
                                        {fileType === constants.MEDIA_TYPE.IMAGE ? (
                                            <img src={fileUrl} alt='Uploaded File' className='max-w-[30%]' />
                                        ) : (
                                            <button>
                                                <svg
                                                    xmlns='http://www.w3.org/2000/svg'
                                                    fill='none'
                                                    viewBox='0 0 24 24'
                                                    strokeWidth={1.5}
                                                    stroke='blue'
                                                    className='w-12'>
                                                    <path
                                                        strokeLinecap='round'
                                                        strokeLinejoin='round'
                                                        d='M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z'
                                                    />
                                                    {fileType === constants.MEDIA_TYPE.VIDEO && (
                                                        <path
                                                            strokeLinecap='round'
                                                            strokeLinejoin='round'
                                                            transform='translate(6.5, 8.25) scale(0.5)'
                                                            d='M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971L6.917 18.307a1.125 1.125 0 01-1.667-.985V5.653z'
                                                        />
                                                    )}
                                                </svg>
                                            </button>
                                        )}
                                        <div className='w-3/5'>
                                            <div className='selectedMediaName'>
                                                {selectedFile?.name}
                                            </div>
                                            <div className='selectedMediaDetails'>
                                                {selectedFile?.size && formatFileSize(selectedFile?.size)}
                                            </div>
                                            {imageLoader && (
                                                <div className='mediaLoaderStyle'>
                                                    <div
                                                        className='mediaUploadProgress'
                                                        style={{ width: `${uploadProgress}%` }}
                                                    />
                                                    <div className='mediaUploadProgressText'>{`${uploadProgress.toFixed(
                                                        0
                                                    )}%`}</div>
                                                </div>
                                            )}
                                        </div>

                                        <svg
                                            fill='gray'
                                            height='10px'
                                            width='10px'
                                            version='1.1'
                                            xmlns='http://www.w3.org/2000/svg'
                                            viewBox='0 0 490 490'
                                            className='absolute top-1.5 right-1.5 cursor-pointer'
                                            xmlSpace='preserve'
                                            onClick={() => removeSelectedFile()}>
                                            <g id='SVGRepo_bgCarrier' strokeWidth='0'></g>
                                            <g
                                                id='SVGRepo_tracerCarrier'
                                                strokeLinecap='round'
                                                strokeLinejoin='round'
                                                stroke='gray'
                                                strokeWidth='0.9800000000000001'></g>
                                            <g id='SVGRepo_iconCarrier'>
                                                <polygon points='456.851,0 245,212.564 33.149,0 0.708,32.337 212.669,245.004 0.708,457.678 33.149,490 245,277.443 456.851,490 489.292,457.678 277.331,245.004 489.292,32.337 '></polygon>{' '}
                                            </g>
                                        </svg>
                                    </div>
                                )}

                                {showError.image && (
                                    <div className='text-red-600 text-sm'>Image above 10MB cannot be uploaded</div>
                                )}
                                {showError.file && (
                                    <div className='text-red-600 text-sm'>File above 50MB cannot be uploaded</div>
                                )}

                                {/* Chat Screen Editor */}
                                <div className='editorBlock'>
                                    <div className='w-full'>
                                        <NewEditor
                                            handleFileChange={handleFileChange}
                                            value={typedMessage}
                                            handleChangeMessage={handleChangeMessage}
                                            userName={selectedChannelState?.channel_name}
                                            setMentionpopup={setMentionpopup}
                                            handleonKeyup={handleonKeyup}
                                            mentionPopup={mentionPopup}
                                        />
                                    </div>
                                    <button
                                        className='submitEditorButton'
                                        disabled={
                                            (!typedMessage.trim().length > 0 || !cleanHTML(typedMessage)) && !fileName.length
                                        }
                                        type='submit'>
                                        <img src={submitIcon} alt="submit icon" style={{ width: 25, height: 25 }} />
                                    </button>
                                </div>
                            </form>
                        )}
                    </div>
                </div>) : (
                <DefaultChatScreen />
            )
            }
            {openPopConfirm && (
                <div className='w-full px-2 py-1'>
                    <ConfirmationModal
                        isOpen={true}
                        onClose={handleClose}
                        title='Confirm Message Delete'
                        message='Are you sure you want to delete this Message?'
                        confirmLabel='Delete'
                        cancelLabel='Cancel'
                        onConfirm={handleDeleteMessage}
                    />
                </div>)}
            <ToastContainer />
            {playVideo && (
                <Video
                    videoUrl={videoUrl}
                    setPlayVideo={setPlayVideo}
                    selectedVideoMessage={selectedVideoMessage}
                    setSelectedVideoMessage={setSelectedVideoMessage}
                    DownloadFile={DownloadFile}
                />
            )}
        </div >
    );
};

ChatScreen.propTypes = {
    isOpen: PropTypes.func?.isRequired,
    fetchChat: PropTypes.func?.isRequired,
    page: PropTypes.number.isRequired,
    total: PropTypes.number,
    setLoading: PropTypes.func.isRequired,
    loading: PropTypes.bool.isRequired,
    userId: PropTypes.string,
};

const mapStateToProps = state => ({
    directChannelState: state.DirectChannelReducer || [],
    favouriteChannelState: state.FavouriteChannelReducer || [],
    selectedChannelState: state.SelectedChannelReducer || {},
    chatHistoryState: state.ChatHistoryReducer || [],
});

export default connect(
    mapStateToProps
)(ChatScreen);

