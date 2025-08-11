import React, { Fragment } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import PropTypes from 'prop-types';

const ConfirmationModal = ({
    isOpen,
    onClose,
    title,
    message,
    confirmLabel,
    cancelLabel,
    onConfirm
}) => {
    return (
        <Transition.Root show={isOpen} as={Fragment}>
            <Dialog as='div' className='dialog-container' onClose={onClose}>
                <Transition.Child
                    as={Fragment}
                    enter='ease-out duration-300'
                    enterFrom='opacity-0'
                    enterTo='opacity-100'
                    leave='ease-in duration-200'
                    leaveFrom='opacity-100'
                    leaveTo='opacity-0'>
                    <div className='overlay' />
                </Transition.Child>

                <div className='modal-container'>
                    <div className='modal-content'>
                        <Transition.Child
                            as={Fragment}
                            enter='ease-out duration-300'
                            enterFrom='opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95'
                            enterTo='opacity-100 translate-y-0 sm:scale-100'
                            leave='ease-in duration-200'
                            leaveFrom='opacity-100 translate-y-0 sm:scale-100'
                            leaveTo='opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95'>
                            <Dialog.Panel className='modal-panel'>
                                <Dialog.Title as='h1' className='dialog-title'>
                                    {title}
                                </Dialog.Title>
                                <div className='dialog-message'>
                                    {message}
                                </div>
                                <div className='modal-buttons'>
                                    <button
                                        type='button'
                                        className='button-confirm'
                                        onClick={onConfirm}>
                                        {confirmLabel}
                                    </button>
                                    <button
                                        type='button'
                                        className='button-cancel'
                                        onClick={onClose}>
                                        {cancelLabel}
                                    </button>
                                </div>
                            </Dialog.Panel>
                        </Transition.Child>
                    </div>
                </div>
            </Dialog>
        </Transition.Root>
    );
};

ConfirmationModal.propTypes = {
    isOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    title: PropTypes.string.isRequired,
    message: PropTypes.string.isRequired,
    confirmLabel: PropTypes.string.isRequired,
    cancelLabel: PropTypes.string.isRequired,
    onConfirm: PropTypes.func.isRequired
};

export default ConfirmationModal;
