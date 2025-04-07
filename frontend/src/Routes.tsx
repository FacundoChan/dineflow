import React, { useState } from "react";
import { Routes, Route } from "react-router-dom";
import SubmitOrder from "./components/SubmitOrder";
import AddItem from "./components/AddItem";
import ItemDisplay from "./components/ItemDisplay";
import PaymentResult from "./pages/PaymentResult";

const MainPage: React.FC = () => {
  const [orderItems, setOrderItems] = useState<
    { id: string; quantity: number }[]
  >([]);

  const handleAddItem = (item: { id: string; quantity: number }) => {
    setOrderItems((prevItems) => [...prevItems, item]);
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="mb-4 text-3xl font-bold">商品订单系统</h1>
      <AddItem onAdd={handleAddItem} />
      <div className="mb-4">
        {orderItems.map((item, index) => (
          <ItemDisplay key={index} item={item} />
        ))}
      </div>
      <SubmitOrder items={orderItems} />
    </div>
  );
};

const AppRoutes: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<MainPage />} />
      <Route path="/success" element={<PaymentResult />} />
    </Routes>
  );
};

export default AppRoutes;
