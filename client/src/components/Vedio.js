import React, { Fragment } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import PropTypes from 'prop-types';
import { ArrowDownIcon } from '@heroicons/react/20/solid';
import { connect } from 'react-redux';
import "../assets/style/components/Vedio.css";

const Video = ({
  videoUrl,
  setPlayVideo,
  selectedChannel,
  selectedVideoMessage,
  setSelectedVideoMessage,
  DownloadFile
}) => {
  const handleCloseModal = () => {
    setPlayVideo(false);
    setSelectedVideoMessage({});
  };

  return (
    <>
      <Transition.Root show={true} as={Fragment}>
        <Dialog as='div' className='dialog-container' onClose={handleCloseModal}>
          <div className='dialog-container-wrapper'>
            <Transition.Child
              as={Fragment}
              enter='ease-out duration-300'
              enterFrom='opacity-0'
              enterTo='opacity-100'
              leave='ease-in duration-200'
              leaveFrom='opacity-100'
              leaveTo='opacity-0'>
              <Dialog.Overlay className='dialog-overlay' />
            </Transition.Child>

            <Transition.Child
              as={Fragment}
              enter='ease-out duration-300'
              enterFrom='opacity-0 scale-95'
              enterTo='opacity-100 scale-100'
              leave='ease-in duration-200'
              leaveFrom='opacity-100 scale-100'
              leaveTo='opacity-0 scale-95'>
              <div className='icon-container'>
                {/* Close Icon and Download Button */}
                <div className='icon-container-wrapper'>
                  <button className='text-white'>
                    <ArrowDownIcon
                      className='h-6 w-6'
                      onClick={() => DownloadFile(selectedVideoMessage)}
                    />
                  </button>
                  <svg
                    xmlns='http://www.w3.org/2000/svg'
                    fill='none'
                    viewBox='0 0 24 24'
                    stroke='currentColor'
                    className='h-6 w-6 text-white cursor-pointer'
                    onClick={handleCloseModal}>
                    <path
                      strokeLinecap='round'
                      strokeLinejoin='round'
                      strokeWidth='2'
                      d='M6 18L18 6M6 6l12 12'
                    />
                  </svg>
                </div>

                <Dialog.Title as='h1' className='dialog-title-vedio'>
                  {selectedVideoMessage?.file_name.substring(
                    selectedVideoMessage?.file_name.lastIndexOf('/') + 1
                  ) || ''}
                </Dialog.Title>
                <Dialog.Description className='dialog-description-vedio'>
                  <span className='text-white'>{selectedVideoMessage?.user?.name} </span>
                  <span>Shared in </span>
                  <span>{selectedChannel?.channel_name}</span>
                </Dialog.Description>

                <div className='vedio-url-container'>
                  <video className='vedio-url-container-wrapper' controls>
                    <source src={videoUrl} type='video/mp4' />
                  </video>
                </div>
              </div>
            </Transition.Child>
          </div>
        </Dialog>
      </Transition.Root>
    </>
  );
};

Video.propTypes = {
  setPlayVideo: PropTypes.func?.isRequired,
  videoUrl: PropTypes.string?.isRequired,
  selectedVideoMessage: PropTypes.object?.isRequired,
  setSelectedVideoMessage: PropTypes.func?.isRequired,
  selectedUser: PropTypes.object?.isRequired,
  DownloadFile: PropTypes.func?.isRequired
};


const mapStateToProps = state => ({
    selectedChannelState: state.SelectedChannelReducer || {},
});

export default connect(
    mapStateToProps
)(Video);

