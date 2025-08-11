import React from 'react';
import routes from "./routes/RouteFile";
import { Route, Routes, BrowserRouter } from 'react-router-dom';
import SomethingwentWrong from './pages/SomethingwentWrong';
import ProtectedRoute from './routes/ProtectedRoutes';

function App() {

  return (
    <BrowserRouter>
      <Routes>
        {routes?.map((route, index) => {
          if (route.private) {
            return (
              <Route
                key={index}
                path='/'
                element={<ProtectedRoute route={route}></ProtectedRoute>}
              >
                <Route path={route.path} element={<route.component />}>
                  {route.children &&
                    route.children.map((childRoute, childIndex) => (
                      <Route
                        key={childIndex}
                        path={childRoute.path}
                        element={<childRoute.component />}
                      />
                    ))}
                </Route>
              </Route>
            );
          }

          return (
            <Route
              key={index}
              path={route.path}
              element={<route.component />}
            />
          );
        })}
        <Route element={<ProtectedRoute />}>
          <Route path='*' element={<SomethingwentWrong />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
