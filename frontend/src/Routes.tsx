// src/Routes.tsx

import React from "react";
import { Routes, Route } from "react-router-dom";
import MainPage from "./pages/MainPage";
import PaymentResult from "./pages/PaymentResult";

const AppRoutes: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<MainPage />} />
      <Route path="/success" element={<PaymentResult />} />
    </Routes>
  );
};

export default AppRoutes;
