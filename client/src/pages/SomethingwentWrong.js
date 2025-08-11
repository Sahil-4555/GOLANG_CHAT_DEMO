import React from 'react';
import { useLocation } from 'react-router-dom';
// import '../assets/style/pages/SomethingwentWrong.css';

const SomethingwentWrong = () => {
  const location = useLocation();
  const errorMessage = location.state && location.state.errorMessage;

  return (
    <section className='something-went-wrong-container'>
      <div className='something-went-wrong-container-wrapper'>
        <div className='error-image '>
          <svg className='error-image-text' fill='none' viewBox='0 0 24 24' stroke='currentColor'>
            <path strokeLinecap='round' strokeLinejoin='round' strokeWidth={2} d='M6 18L18 6M6 6l12 12' />
          </svg>
        </div>
        <h1 className='oops-text-header'>Oops!</h1>
        <div className='oops-container'>
          <h3 className='error-as-header'>Looks like we're lost</h3>
          {errorMessage && <p className='mb-4'>{errorMessage}</p>}
          <a href='/login' className='route-style'>
            Go to Home
          </a>
        </div>
      </div>
    </section>
  );
}

export default SomethingwentWrong;    