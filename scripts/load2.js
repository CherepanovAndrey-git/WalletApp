// old k6 tested for service load


// import http from 'k6/http';
// import { check, sleep } from 'k6';
//
// export let options = {
//   stages: [
//
//     { duration: '20s', target: 2000 },
//     { duration: '1m', target: 2000 },
//     { duration: '40s', target: 0 },
//   ],
//   thresholds: {
//     http_req_failed: ['rate<0.01'],
//     http_req_duration: ['p(95)<500'],
//   },
// };
//
// export default function () {
//   const walletId = "WalletId"; //Generate first
//   const apiKey = "ApiKey";  //Generate first
//   const baseUrl = "http://localhost:8080/v1";
//
//   const headers = {
//     'Content-Type': 'application/json',
//     'ApiKey': apiKey,
//   };
//
//
//   const depositPayload = JSON.stringify({
//     walletId: walletId,
//     operationType: "DEPOSIT",
//     amount: 100,
//   });
//
//   let res = http.post(`${baseUrl}/wallet`, depositPayload, { headers: headers });
//   check(res, {
//     'deposit status is 200': (r) => r.status === 200,
//     'deposit balance updated': (r) => JSON.parse(r.body).balance !== undefined,
//   });
//
//
//   const withdrawPayload = JSON.stringify({
//     walletId: walletId,
//     operationType: "WITHDRAW",
//     amount: 50,
//   });
//
//   res = http.post(`${baseUrl}/wallet`, withdrawPayload, { headers: headers });
//   check(res, {
//     'withdraw status is 200': (r) => r.status === 200,
//     'withdraw balance updated': (r) => JSON.parse(r.body).balance !== undefined,
//   });
//
//
//   res = http.get(`${baseUrl}/wallets/${walletId}/balance`, { headers: headers });
//   check(res, {
//     'balance status is 200': (r) => r.status === 200,
//     'balance is valid': (r) => JSON.parse(r.body).balance !== undefined,
//   });
//
//   sleep(1);
// }
