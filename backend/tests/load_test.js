import http from "k6/http";
import { sleep } from "k6";

export const options = {
  stages: [
    { duration: "1m", target: 100 },
    { duration: "1m", target: 100 },
    { duration: "1m", target: 0 },
  ],
  threshold: {
    http_req_duration: ["p(99)<100"],
  },
};

export default function () {
  const url = "http://127.0.0.1:54431/api/v1/courses/paginated";
  const payload = JSON.stringify({
    
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const response = http.post(url, payload, params)

  if (response.status !== 200) {
    console.error(`Non-200 response: ${response.error}`);
  }

  sleep(1);
}
