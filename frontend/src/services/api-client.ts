import applyInterceptors from "@/middleware/axiosMiddleware";
import axios from "axios";

const apiClient = axios.create({
  baseURL: "http://127.0.0.1:8282",
});

applyInterceptors(apiClient);

export default apiClient;
