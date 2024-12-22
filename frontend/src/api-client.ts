import axios from "axios";

const PORT = "8181"
const SERVER_URL = `http://localhost:${PORT}/api`

export default axios.create({
    baseURL: SERVER_URL,
    headers: {
      'Content-Type': 'application/json',
    },
  });
