import React, { useEffect, useRef, useState } from 'react';
import PropTypes from 'prop-types';
// import * as Emoji from 'quill-emoji';
import 'quill-emoji/dist/quill-emoji.css'; // Import the emoji styles;
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css'; // ES6
import { connect } from 'react-redux';
import paperclip from '../assets/icons/paperclip.svg'

const formats = [
  'header',
  'bold',
  'italic',
  'underline',
  'align',
  'strike',
  'script',
  'blockquote',
  'background',
  'list',
  'bullet',
  'indent',
  'link',
  'image',
  'color',
  'code-block',
  'emoji',
  'paperclip',
  'mention'
];

const toolBarOptions = [
  ['bold', 'italic', 'underline'],
  [{ list: 'ordered' }, { list: 'bullet' }, { indent: '-1' }, { indent: '+1' }],
  ['code-block'],
  [('clean', 'paperclip')],
  ['emoji']
];

const modules = {
  toolbar: {
    container: toolBarOptions
  },
  clipboard: {
    matchVisual: false
  },
  'emoji-toolbar': true,
  'emoji-textarea': false
};

export const NewEditor = (props) => {
  const {
    handleFileChange,
    value,
    handleChangeMessage,
    userName,
    handleonKeyup,
    setMentionpopup,
    mentionPopup,
    selectedChannelState,
    usersState
  } = props;

  const quillRef = useRef(null);
  const [key, setKey] = useState(0);

  useEffect(() => {
    setKey((prevKey) => prevKey + 1);
  }, []);

  modules['emoji-shortname'] = {
    fuse: {
      shouldSort: true,
      threshold: 0.1,
      location: 1,
      distance: 100,
      maxPatternLength: 32,
      minMatchCharLength: 2,  
      keys: ['shortname']
    },
    onOpen: function () {
      setMentionpopup(true);
    },
    onClose: function () {
      setTimeout(() => {
        setMentionpopup(false);
      }, 100);
    }
  };

  modules.mention = {
    allowedChars: /^[A-Za-z\sÅÄÖåäö]*$/,
    mentionDenotationChars: ['@', '#'],
    source: function (searchTerm, renderList, mentionChar) {
      let values;
      const allUsers = usersState?.map((item) => ({ value: item?.name, id: item._id }));
      if (mentionChar === '@') {
        values = selectedChannelState?.users?.map((item) => ({ value: item?.name, id: item._id }));
        if (searchTerm.length === 0) {
          renderList(values, searchTerm);
        } else {
          const matches = [];
          for (let i = 0; i < allUsers.length; i++)
            if (~allUsers[i].value.toLowerCase().indexOf(searchTerm.toLowerCase()))
              matches.push(allUsers[i]);
          renderList(matches, searchTerm);
        }
      }
    },
    onOpen: function () {
      setMentionpopup(true);
    },
    onClose: function () {
      setTimeout(() => {
        setMentionpopup(false);
      }, 100);
    }
  };

  modules.keyboard = {
    bindings: {
      handleEnter: {
        key: 13,
        handler: function (range, context) {
          if (context.event) {
            if (mentionPopup) {
              setMentionpopup(false);
              return;
            }
            context.event.preventDefault();
          }
        }
      }
    }
  };

  modules.toolbar.handlers = {
    paperclip: handleFileChange
  };

  useEffect(() => {
    const button = document.getElementsByClassName('ql-paperclip')[0];
    if (button) button.innerHTML = paperclip;
  }, []);

  useEffect(() => {
    if (quillRef?.current) {
      const quill = quillRef.current.getEditor();
      quill.root.dataset.placeholder = `Write to ${selectedChannelState?.channel_name}`;
      quill.focus();
    }
  }, [userName, quillRef]);

  return (
    <div className='text-editor' key={key}>
      <ReactQuill
        theme='snow'
        value={value}
        onChange={(e) => handleChangeMessage(e)}
        placeholder={`Write to ${selectedChannelState?.channel_name}`}
        modules={modules}
        formats={formats}
        onKeyDown={(e) => {
          handleonKeyup && handleonKeyup(e);
          // handleEditor(e);
        }}
        ref={quillRef}
      />
    </div>
  );
};

NewEditor.propTypes = {
  handleFileChange: PropTypes.func?.isRequired,
  value: PropTypes.string?.isRequired,
  handleChangeMessage: PropTypes.func?.isRequired,
  handleonKeyup: PropTypes.func,
  userName: PropTypes.string?.isRequired,
  setMentionpopup: PropTypes.func?.isRequired,
  mentionPopup: PropTypes.bool?.isRequired,
};

const mapStateToProps = state => ({
  selectedChannelState: state.SelectedChannelReducer || {},
  usersState: state.UsersReducer || [],
});

export default connect(
  mapStateToProps
)(NewEditor);
