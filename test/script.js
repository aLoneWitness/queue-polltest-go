import http from 'k6/http'
import {check} from 'k6'
import {uuidv4} from 'https://jslib.k6.io/k6-utils/1.1.0/index.js';

export let options = {
    stages: [
        { duration: '10s', target: 5000 },
        { duration: '50s', target: 5000 },
    ],
};

export function setup() {
    // 2. setup code, you can pass data to VU and teardown

    return {}
}

export default function (data) {
    // 3. VU code
    const baseUrl = 'http://localhost:8080/';

    let uuid = uuidv4()

    const params = {
        headers: { 'LB_HEADER_AFFINITY': uuid },
    };

    const response = http.get(baseUrl + 'queue', params)
    check(response, {
        'response code was 200': (response) => response.status === 200,
    });
}

export function teardown(data) {

}