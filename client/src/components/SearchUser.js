import React, { useEffect, useRef, useState } from "react";
import PropTypes from "prop-types";
import * as constants from "../utils/constant/Constant";
import "../assets/style/components/SearchUser.css";

function SearchUser(props) {
    const {
        searchedChannel,
        onSelectUser,
        searchedvalue,
        setSearchedvalue,
        isOpen,
        searchisOpen,
        onClose,
    } = props;

    const focusRef = useRef(null);
    const [focusedIndex, setFocusedIndex] = useState(0);
    const optionsContainerRef = useRef(null);
    const optionsListRef = useRef(null);
    const outsideRef = useRef(null);
    const insideRef = useRef(null);

    const listFocus = (e) => {
        switch (e.key) {
            case "ArrowDown":
                e.preventDefault();
                setFocusedIndex((prevIndex) =>
                    prevIndex < searchedChannel.length - 1 ? prevIndex + 1 : 0
                );
                break;
            case "ArrowUp":
                e.preventDefault();
                setFocusedIndex((prevIndex) =>
                    prevIndex > 0 ? prevIndex - 1 : searchedChannel.length - 1
                );
                break;
            case "Enter":
                if (focusedIndex !== -1) {
                    handleSelect(searchedChannel[focusedIndex]);
                    onClose(false)
                }
                break;
            default:
                break;
        }
    };

    const handleKeyDown = (e) => {
        if (e.ctrlKey && e.key === "k") {
            e.preventDefault();
            isOpen(true);
            queueMicrotask(() => {
                focusRef.current.focus();
                setSearchedvalue("");
                setFocusedIndex(0);
            });
        }

        if (e.key === "Escape") {
            onClose(false);
        }
    };

    const handleSelect = (option) => {
        if (!option) return;
        onSelectUser(option, constants.SELECTION_TYPES.FROM_SEARCH);
    };

    const handleClickOutside = (e) => {
        if (
            (outsideRef.current && !insideRef.current) ||
            !insideRef.current.contains(e.target)
        ) {
            onClose(false);
        }
    };

    useEffect(() => {
        if (
            focusedIndex !== -1 &&
            optionsListRef.current &&
            optionsContainerRef.current
        ) {
            const listItem = optionsListRef.current.children[focusedIndex];
            if (listItem) {
                listItem.scrollIntoView({
                    block: "nearest",
                    inline: "nearest",
                    behavior: "smooth",
                });
            }
        }
    }, [focusedIndex]);

    useEffect(() => {
        document.addEventListener("keydown", handleKeyDown);
        outsideRef.current?.addEventListener("click", handleClickOutside);

        return () => {
            document.removeEventListener("keydown", handleKeyDown);
            outsideRef.current?.removeEventListener("click", handleClickOutside);
        };
    }, []);

    return (
        <div ref={outsideRef} className={`container ${searchisOpen ? 'visible' : 'hidden'}`}>
            <div className="dialog" ref={insideRef}>
                <div className="dialog-container-search-user ">
                    <h5 className="dialog-title-search-user">Find Channels</h5>
                    <button onClick={onClose}>
                        <svg
                            fill="#000000"
                            height="20px"
                            width="20px"
                            version="1.1"
                            id="Capa_1"
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 490 490"
                            xmlSpace="preserve"
                        >
                            <g id="SVGRepo_bgCarrier" strokeWidth="0"></g>
                            <g
                                id="SVGRepo_tracerCarrier"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                stroke="#CCCCCC"
                                strokeWidth="0.9800000000000001"
                            ></g>
                            <g id="SVGRepo_iconCarrier">
                                {" "}
                                <polygon points="456.851,0 245,212.564 33.149,0 0.708,32.337 212.669,245.004 0.708,457.678 33.149,490 245,277.443 456.851,490 489.292,457.678 277.331,245.004 489.292,32.337 "></polygon>{" "}
                            </g>
                        </svg>
                    </button>
                </div>
                <div>
                    <input
                        type="search"
                        ref={focusRef}
                        value={searchedvalue}
                        onKeyDown={(e) => listFocus(e)}
                        onChange={(e) => setSearchedvalue(e.target.value)}
                        className="search-input-box"
                        placeholder="Search your Channel"
                    />
                </div>
                <div className="p-2" ref={optionsContainerRef}>
                    <div>
                        <div>
                            {searchedChannel?.length > 0 ? (
                                <ul role="list" className="options-list" ref={optionsListRef}>
                                    {searchedChannel?.map((user, index) => {
                                        return (
                                            <li
                                                onClick={() => {
                                                    onSelectUser(
                                                        user,
                                                        constants.SELECTION_TYPES.FROM_SEARCH,
                                                    );
                                                    onClose(false);
                                                }}
                                                className={`selected-option-search-user ${focusedIndex === index ? "bg-blue-200" : ""
                                                    }`}
                                                key={index}
                                            >
                                                {user.channel_type === constants.CHANNEL_TYPES.DIRECT_CHANNEL ?
                                                    <div className="options-avatar">
                                                        <div className="options-avatar-name">
                                                            {user?.channel_name?.charAt(0).toUpperCase()}
                                                        </div>
                                                        <div className="options-channel-name-container">
                                                            <p className="options-channel-name">
                                                                {user?.channel_name}
                                                            </p>
                                                        </div>
                                                    </div> :
                                                    <div className="options-avatar">
                                                        <div className="options-avatar-name">
                                                            <svg
                                                                xmlns='http://www.w3.org/2000/svg'
                                                                fill='none'
                                                                viewBox='0 0 24 24'
                                                                strokeWidth={1.5}
                                                                stroke='currentColor'
                                                                className={'w-5 h-5'}>
                                                                <path
                                                                    strokeLinecap='round'
                                                                    strokeLinejoin='round'
                                                                    d='M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z'
                                                                />
                                                            </svg>
                                                        </div>
                                                        <div className="options-channel-name-container">
                                                            <p className="options-channel-name">
                                                                {user?.channel_name}
                                                            </p>
                                                        </div>
                                                    </div>
                                                }
                                            </li>
                                        );
                                    })}
                                </ul>
                            ) : (
                                <div className="message-no-result-found">
                                    No Result Found
                                </div>
                            )}
                        </div>

                    </div>
                </div>
            </div>
        </div>
    );
}

SearchUser.propTypes = {
    searchedChannel: PropTypes.any.isRequired,
    onSelectUser: PropTypes.func.isRequired,
    searchedvalue: PropTypes.any.isRequired,
    setSearchedvalue: PropTypes.func.isRequired,
    isOpen: PropTypes.func.isRequired,
    searchisOpen: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
};

export default SearchUser;