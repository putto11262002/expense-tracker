import axios from "axios";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:3001",
    withCredentials: true
})

api.interceptors.request.use(function(config) {
    // if auth token if not present in the header try check local storage
    if (!config.headers.Authorization){
        const token = localStorage.getItem("token")
        if(token){
            config.headers.Authorization = `Bearer ${token}` 
        }
    }
    return config
}, function(error){
    return Promise.reject(error)
})



export default api;