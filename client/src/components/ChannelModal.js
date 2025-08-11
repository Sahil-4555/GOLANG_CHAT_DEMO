import { Fragment, useState } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { RadioGroup } from '@headlessui/react';
import { CheckCircleIcon } from '@heroicons/react/20/solid';
import api from '../utils/api';
import PropTypes from 'prop-types';
import * as handlers from '../redux/actions/actions'
import * as constants from '../utils/constant/Constant';
import { useDispatch } from 'react-redux';
import "../assets/style/components/ChannelModal.css";

const INPUT_VALUES = {
    NAME: 'name',
    DESCRIPTION: 'description'
};

function classNames(...classes) {
    return classes.filter(Boolean).join(' ');
}

function ChannelModal({ handleChangeModal }) {
    const [selectedChannelType, setSelectedChannelType] = useState(constants.CHANNEL_OPTIONS[1]);
    const [channelName, setChannelName] = useState('');
    const [channelDescription, setChannelDescription] = useState('');
    const [errors, setErrors] = useState({});
    const dispatch = useDispatch();

    const handleChangeValues = (e, type) => {
        const value = e.target.value;
        if (type === INPUT_VALUES.NAME) return setChannelName(value);
        setChannelDescription(value);
    };

    const handleSubmitNewChannel = async () => {
        try {
            if (channelDescription === "") {
                const payload = {
                    channel_name: channelName,
                };
                const data = await api.post('/v1/channel/createGroup', payload);
                if (data?.data?.meta?.code === constants.API_STATUS.FAILURE_CODE) {
                    setErrors({ warning: data?.data?.meta?.message });
                    return;
                }
            } else {
                const payload = {
                    channel_name: channelName,
                    description: channelDescription,
                };
                const data = await api.post('/v1/channel/createGroup', payload);
                if (data?.data?.meta?.code === constants.API_STATUS.FAILURE_CODE) {
                    setErrors({ warning: data?.data?.meta?.message });
                    return;
                }
            }

            handleChangeModal();
            dispatch(handlers.getGroupChannels())
        } catch (err) {
            console.log('Error : ', err);
        }
    };

    return (
        <>
            <Transition.Root show={true} as={Fragment}>
                <Dialog
                    as='div'
                    className='channelModal'
                    onClose={handleChangeModal}
                    style={{ maxWidth: '80rem' }}>
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

                    <div className='dialogContainer'>
                        <div className='dialogContent'>
                            <Transition.Child
                                as={Fragment}
                                enter='ease-out duration-300'
                                enterFrom='opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95'
                                enterTo='opacity-100 translate-y-0 sm:scale-100'
                                leave='ease-in duration-200'
                                leaveFrom='opacity-100 translate-y-0 sm:scale-100'
                                leaveTo='opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95'>
                                <Dialog.Panel className='dialogPanel'>
                                    <Dialog.Title as='h1' className='text-bold font-semibold leading-2 text-gray-900'>
                                        Create a new channel
                                    </Dialog.Title>

                                    <div className='mt-5'>
                                        <div className='relative'>
                                            <label
                                                htmlFor='name'
                                                className='inputLabel'>
                                                channel name
                                            </label>
                                            <input
                                                type='text'
                                                name='name'
                                                autoComplete="off"
                                                id='name'
                                                maxLength={30}
                                                className='inputField'
                                                placeholder='Enter a name of your channel'
                                                onChange={(e) => {
                                                    handleChangeValues(e, INPUT_VALUES.NAME);
                                                    setErrors({});
                                                }}
                                            />
                                            <span className={`inputError ${errors.warning ? '' : ''} `}>
                                                {errors.warning}
                                            </span>
                                        </div>
                                        <RadioGroup value={selectedChannelType} onChange={setSelectedChannelType}>
                                            <div className='radioOptions'>
                                                {constants.CHANNEL_OPTIONS.map((channel) => (
                                                    <RadioGroup.Option
                                                        key={channel.id}
                                                        value={channel}
                                                        className={({ active }) =>
                                                            classNames(
                                                                active ? 'activeBorder' : 'inactiveBorder',
                                                                'radioOption'
                                                            )
                                                        }>
                                                        {({ checked, active }) => (
                                                            <>
                                                                <span className='option-content'>
                                                                    <>
                                                                        <RadioGroup.Label
                                                                            as='span'
                                                                            className='option-label'>
                                                                            {channel.title}
                                                                        </RadioGroup.Label>
                                                                        <RadioGroup.Description
                                                                            as='span'
                                                                            className='option-description'>
                                                                            {channel.description}
                                                                        </RadioGroup.Description>
                                                                    </>
                                                                </span>
                                                                <CheckCircleIcon
                                                                    className={classNames(
                                                                        !checked ? 'invisible' : '',
                                                                        'check-icon'
                                                                    )}
                                                                    aria-hidden='true'
                                                                />
                                                                <span
                                                                    className={classNames(
                                                                        active ? 'border-blue' : 'border-2',
                                                                        checked ? 'border-blue-600' : 'border-transparent',
                                                                        'check-option-border'
                                                                    )}
                                                                    aria-hidden='true'
                                                                />
                                                            </>
                                                        )}
                                                    </RadioGroup.Option>
                                                ))}
                                            </div>
                                        </RadioGroup>
                                        <div className='mt-4'>
                                            <textarea
                                                rows={4}
                                                name='description'
                                                autoComplete="off"
                                                id='description'
                                                className='textarea-input-field'
                                                defaultValue={''}
                                                placeholder='Enter a purpose for this channel (optional)'
                                                onChange={(e) => handleChangeValues(e, INPUT_VALUES.DESCRIPTION)}
                                            />
                                            <label
                                                htmlFor='description'
                                                className='description-label'>
                                                This will be displayed when browsing for channels.
                                            </label>
                                        </div>
                                    </div>

                                    <div className='dialogFooter'>
                                        <button
                                            type='button'
                                            className='cancel-button'
                                            onClick={() => handleChangeModal(false)}>
                                            Cancel
                                        </button>
                                        <button
                                            type='button'
                                            className='create-button'
                                            onClick={() => handleSubmitNewChannel()}
                                            disabled={!channelName.trim()}>
                                            Create channel
                                        </button>
                                    </div>
                                </Dialog.Panel>
                            </Transition.Child>
                        </div>
                    </div>
                </Dialog>
            </Transition.Root>
        </>
    );
}

ChannelModal.propTypes = {
    handleChangeModal: PropTypes.func.isRequired,
};

export default ChannelModal;
