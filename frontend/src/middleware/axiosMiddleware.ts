import { AxiosInstance } from "axios";

const applyInterceptors = (apiClient: AxiosInstance) => {
  apiClient.interceptors.request.use(
    (config) => {
      // Add any request interceptors here
      return config;
    },
    (error) => {
      // Handle request error
      return Promise.reject(error);
    },
  );
};

export default applyInterceptors;
