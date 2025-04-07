import React, { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import axios from "axios";

const PaymentResult: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState("");
  const [hasAlerted, setHasAlerted] = useState(false);

  const customerID = searchParams.get("customerID");
  const orderID = searchParams.get("orderID");

  useEffect(() => {
    const fetchOrderStatus = async () => {
      if (!customerID || !orderID) return;

      try {
        const response = await axios.get(
          `http://127.0.0.1:8282/api/customer/${customerID}/orders/${orderID}`,
        );
        console.log("response.data", response.data);
        const order = response.data.data.Order;
        setStatus(order.status);
      } catch (error) {
        console.error("获取订单状态失败", error);
        setStatus("error");
      }
    };

    fetchOrderStatus();
  }, [customerID, orderID]);

  useEffect(() => {
    if (status !== "paid") return;

    const interval = setInterval(async () => {
      try {
        const response = await axios.get(
          `http://127.0.0.1:8282/api/customer/${customerID}/orders/${orderID}`,
        );
        const order = response.data.data.Order;
        if (order.status === "ready" && !hasAlerted) {
          alert(`订单：${order.id} 已完成`);
          setHasAlerted(true);
          setStatus("ready");
        }
      } catch (error) {
        console.error("轮询失败", error);
      }
    }, 5000);
    return () => clearInterval(interval); // 清理轮询
  }, [status, customerID, orderID, hasAlerted]);
  
  if (status === "paid") return <h1>✅ 支付成功！</h1>;
  if (status === "waiting_for_payment") return <h1>⌛ 等待支付中...</h1>;
  if (status === "error") return <h1>❌ 查询失败</h1>;
  if (status === "ready") return <h1>🎉 订单已完成！</h1>;

  return <h1>🔄 查询中...</h1>;
};

export default PaymentResult;
