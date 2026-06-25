import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAppSelector } from '../store';
import MainLayout from '../layouts/MainLayout';
import LoginPage from '../pages/login';
import ProjectsPage from '../pages/projects';
import BoardPage from '../pages/board';
import IssueDetailPage from '../pages/issue';

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAppSelector((state) => state.auth);
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
};

const AppRouter: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <MainLayout>
                <Navigate to="/projects" replace />
              </MainLayout>
            </PrivateRoute>
          }
        />
        <Route
          path="/projects"
          element={
            <PrivateRoute>
              <MainLayout>
                <ProjectsPage />
              </MainLayout>
            </PrivateRoute>
          }
        />
        <Route
          path="/board"
          element={
            <PrivateRoute>
              <MainLayout>
                <BoardPage />
              </MainLayout>
            </PrivateRoute>
          }
        />
        <Route
          path="/issues/:issueKey"
          element={
            <PrivateRoute>
              <MainLayout>
                <IssueDetailPage />
              </MainLayout>
            </PrivateRoute>
          }
        />
        <Route path="*" element={<Navigate to="/projects" replace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default AppRouter;
