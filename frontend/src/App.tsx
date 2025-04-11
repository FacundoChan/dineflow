import { BrowserRouter } from "react-router-dom";
import AppRoutes from "./Routes";
import { Toaster } from "@/components/ui/sonner";
import "./App.css";

function App() {
  return (
    <BrowserRouter>
      <AppRoutes />
      <Toaster
        richColors
        position="bottom-center"
        toastOptions={{
          className: "text-center",
          descriptionClassName: "text-center",
        }}
      />
    </BrowserRouter>
  );
}

export default App;
