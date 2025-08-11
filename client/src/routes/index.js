import React from 'react';
import { Outlet } from 'react-router-dom';
import PropTypes from 'prop-types';

const Layout = ({ children }) => {
  return (
    <div>
      {children ? (
        children
      ) : (
        <Outlet />
      )}
    </div>
  );
};

export default Layout;

Layout.propTypes = {
  children: PropTypes.node,
};
