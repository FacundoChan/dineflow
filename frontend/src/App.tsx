import { useState } from "react";
// import reactLogo from './assets/react.svg'
// import viteLogo from '/vite.svg'
import "./App.css";
import ItemDisplay from "./components/ItemDisplay";
import SubmitOrder from "./components/SubmitOrder";
import AddItem from "./components/AddItem";

function App() {
  const [orderItems, setOrderItems] = useState<
    { ID: string; Quantity: number }[]
  >([]);

  const handleAddItem = (item: { ID: string; Quantity: number }) => {
    // 每次添加商品后将新商品追加到订单列表中
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
}

// function App() {
//   const [count, setCount] = useState(0)

//   return (
//     <>
//       <div>
//         <a href="https://vite.dev" target="_blank">
//           <img src={viteLogo} className="logo" alt="Vite logo" />
//         </a>
//         <a href="https://react.dev" target="_blank">
//           <img src={reactLogo} className="logo react" alt="React logo" />
//         </a>
//       </div>
//       <h1>Vite + React</h1>
//       <div className="card">
//         <button onClick={() => setCount((count) => count + 1)}>
//           count is {count}
//         </button>
//         <p>
//           Edit <code>src/App.tsx</code> and save to test HMR
//         </p>
//       </div>
//       <p className="read-the-docs">
//         Click on the Vite and React logos to learn more
//       </p>
//     </>
//   )
// }

export default App;
