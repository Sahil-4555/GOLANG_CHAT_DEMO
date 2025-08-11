import { Navigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import Layout from '.';


const checkTokenValid = () => {
  let token = localStorage.getItem('token');

  return token || false;
};

const ProtectedRoute = ({ children }) => {

  if (!checkTokenValid()) {
    return <Navigate to='/login' replace />;
  }

  return <Layout>{children}</Layout>;

};

export default ProtectedRoute;

ProtectedRoute.propTypes = {
  component: PropTypes.object,
  children: PropTypes.node,
};
