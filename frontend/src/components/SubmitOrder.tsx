// src/components/SubmitOrder.tsx
import React, { useEffect, useRef, useState } from "react";
import apiClient from "@/services/api-client";

interface SubmitOrderProps {
  items: { id: string; quantity: number }[];
}

const SubmitOrder: React.FC<SubmitOrderProps> = ({ items }) => {
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [paymentLink, setPaymentLink] = useState("");
  const [isPolling, setIsPolling] = useState(false);
  const pollingIntervalRef = useRef<NodeJS.Timeout | null>(null);

  const generateCustomerID = () => {
    // 生成一个8位随机字符串作为customerID
    return Math.random().toString(36).substring(2, 10);
  };

  const pollOrderStatus = (customerID: string, orderID: string) => {
    pollingIntervalRef.current = setInterval(async () => {
      let order;
      try {
        console.log(
          "GET ",
          `/api/customer/${customerID}/orders/${orderID}`,
        );
        const response = await apiClient.get(
          `/api/customer/${customerID}/orders/${orderID}`,
        );

        console.log(JSON.stringify(response.data));
        order = response.data.data.Order;

        console.log(
          "showPaymentModal:",
          showPaymentModal,
          " Order status: ",
          order.status,
        );
        if (!showPaymentModal && order.status === "waiting_for_payment") {
          setPaymentLink(order.payment_link);
          console.log("link: ", order.payment_link);
          setShowPaymentModal(true);
          setIsPolling(false);
        }
        if (order.status === "paid") {
          alert("支付成功");
          if (pollingIntervalRef.current !== null) {
            clearInterval(pollingIntervalRef.current);
          }
        }
      } catch (error) {
        console.error("轮询获取订单信息失败", error);
      }
    }, 5000);
  };

  const handleSubmit = async () => {
    const customerID = generateCustomerID();
    const postData = {
      customer_id: customerID,
      items: items,
    };

    console.log(
      "POST address:",
      `/api/customer/${customerID}/orders`,
    );
    console.log("提交订单", JSON.stringify(postData));

    try {
      const response = await apiClient.post(
        `/api/customer/${customerID}/orders`,
        postData,
      );
      console.log("Submit order successfully: ", response.data);

      setIsPolling(true);
      console.log("response.data", response.data);
      pollOrderStatus(
        response.data.data.customer_id,
        response.data.data.order_id,
      );
    } catch (error) {
      console.error("Order Submit Error", error);
    }
  };

  const handlePayment = () => {
    // 示例：点击支付按钮后跳转到支付页面
    window.location.href = paymentLink;
  };

  useEffect(() => {
    return () => {
      if (pollingIntervalRef.current !== null) {
        clearInterval(pollingIntervalRef.current);
      }
    };
  }, []);

  return (
    <div>
      <button
        className="rounded bg-green-500 p-2 text-white hover:bg-green-600"
        onClick={handleSubmit}
      >
        提交订单
      </button>
      {isPolling && (
        <div className="mt-4 flex items-center justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
        </div>
      )}
      {showPaymentModal && (
        <div className="bg-opacity-50 fixed inset-0 flex items-center justify-center bg-gray-400">
          <div className="rounded-lg bg-white px-6 py-8 shadow-xl ring ring-gray-900/5 dark:bg-gray-800">
            <h2 className="mt-5 text-base font-medium tracking-tight text-gray-900 dark:text-white">
              支付窗口
            </h2>
            <p className="mt-2 mb-4 text-sm text-gray-500 dark:text-gray-400">
              请点击下面按钮进行支付
            </p>
            <button
              className="dark:text-white-400 rounded bg-blue-500 p-2 text-white hover:bg-sky-900 dark:bg-sky-800"
              onClick={handlePayment}
            >
              立即支付
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default SubmitOrder;
