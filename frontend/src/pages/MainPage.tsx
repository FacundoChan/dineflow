import React, { useEffect, useRef, useState } from "react";
import ProductList from "../components/ProductList";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import apiClient from "@/services/api-client";

type Product = {
  id: string;
  name: string;
  quantity: number;
  price: number;
};

type OrderItem = {
  id: string;
  quantity: number;
};

const MainPage: React.FC = () => {
  const [orderItems, setOrderItems] = useState<OrderItem[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [paymentLink, setPaymentLink] = useState("");
  const pollingIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const getProductById = (id: string) => products.find((p) => p.id === id);

  const getTotalAmount = () => {
    return orderItems.reduce((sum, item) => {
      const product = getProductById(item.id);
      return product ? sum + product.price * item.quantity : sum;
    }, 0);
  };

  const generateCustomerID = () => {
    return Math.random().toString(36).substring(2, 10);
  };

  const pollOrderStatus = (customerID: string, orderID: string) => {
    pollingIntervalRef.current = setInterval(async () => {
      try {
        const response = await apiClient.get(
          `/api/customer/${customerID}/orders/${orderID}`,
        );
        const order = response.data.data.Order;

        if (!showPaymentModal && order.status === "waiting_for_payment") {
          setPaymentLink(order.payment_link);
          setShowPaymentModal(true);
        }

        if (order.status === "paid") {
          toast("支付成功", { description: "订单已完成" });
          if (pollingIntervalRef.current)
            clearInterval(pollingIntervalRef.current);
        }
      } catch (err) {
        console.error("轮询失败", err);
      }
    }, 500);
  };

  const handleSubmitOrder = async () => {
    const customerID = generateCustomerID();
    const postData = { customer_id: customerID, items: orderItems };

    setSubmitting(true);
    toast("订单提交中...", { description: "正在生成订单，请稍候..." });

    try {
      const response = await apiClient.post(
        `/api/customer/${customerID}/orders`,
        postData,
      );
      toast("订单已提交", { description: "等待付款中" });

      const data = response.data.data;
      console.log("response.data", data);

      if (data.Order?.status === "waiting_for_payment") {
        setPaymentLink(data.Order.payment_link);
        setShowPaymentModal(true);
        return;
      }

      pollOrderStatus(data.customer_id, data.order_id);
    } catch (err) {
      toast("订单提交失败", {
        description: "请检查网络或稍后重试",
      });
    } finally {
      setSubmitting(false);
    }
  };

  const handlePayment = () => {
    window.location.href = paymentLink;
  };

  useEffect(() => {
    return () => {
      if (pollingIntervalRef.current) clearInterval(pollingIntervalRef.current);
    };
  }, []);

  return (
    <div className="container mx-auto p-4">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-3xl font-bold">商品订单系统</h1>
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="outline">查看购物车</Button>
          </SheetTrigger>
          <SheetContent className="p-3">
            <SheetHeader>
              <SheetTitle>订单清单</SheetTitle>
              <SheetDescription>目前已选择的商品如下：</SheetDescription>
            </SheetHeader>
            <div className="mt-2 space-y-3 p-2.5">
              {orderItems.length === 0 ? (
                <p className="text-muted-foreground flex justify-center">
                  当前未选择任何商品
                </p>
              ) : (
                orderItems.map((item) => {
                  const product = getProductById(item.id);
                  const subtotal = product ? product.price * item.quantity : 0;

                  return (
                    <div key={item.id} className="flex flex-col pb-1">
                      <div className="flex justify-between text-sm font-bold">
                        <span>{product?.name || item.id}</span>
                        <span>x {item.quantity}</span>
                      </div>
                      <div className="flex justify-between text-xs text-gray-500">
                        <span>HK${product?.price.toFixed(2)} / 件</span>
                        <span>小计: HK${subtotal.toFixed(2)}</span>
                      </div>
                    </div>
                  );
                })
              )}
              <hr className="my-2" />
              <div className="flex justify-between font-bold">
                <span>总数量</span>
                <span>
                  {orderItems.reduce((sum, item) => sum + item.quantity, 0)} 件
                </span>
              </div>
              <div className="flex justify-between font-bold text-slate-600">
                <span>总金额</span>
                <span>HK${getTotalAmount().toFixed(2)}</span>
              </div>
              <Button
                className="mt-4 w-full"
                disabled={submitting || orderItems.length === 0}
                onClick={handleSubmitOrder}
              >
                {submitting ? "提交中..." : "提交订单"}
              </Button>
            </div>
          </SheetContent>
        </Sheet>
      </div>
      <ProductList
        onChange={setOrderItems}
        products={products}
        setProducts={setProducts}
      />
      {showPaymentModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div className="mt-3 rounded-lg bg-white p-6 shadow-lg dark:bg-gray-900">
            <h2 className="mb-3 text-lg font-semibold">支付窗口</h2>
            <p className="mb-4 text-sm text-gray-500 dark:text-gray-400">
              订单已创建，请点击下方按钮前往支付页面。
            </p>
            <Button className="mb-2 w-auto" onClick={handlePayment}>
              立即支付
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

export default MainPage;
