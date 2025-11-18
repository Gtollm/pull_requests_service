import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const NUM_USERS = parseInt(__ENV.NUM_USERS || '10');

export const options = {
    vus: 20,
    duration: '30s',
};

export function setup() {
    const members = [];
    for (let i = 1; i <= NUM_USERS; i++) {
        const paddedId = String(i).padStart(2, '0');
        members.push({
            user_id: `550e8400-e29b-41d4-a716-4466554400${paddedId}`,
            username: `User${i}`,
            is_active: true
        });
    }

    const teamName = `team-k6-${Date.now()}`;

    const teamRes = http.post(
        `${BASE_URL}/team/add`,
        JSON.stringify({ team_name: teamName, members: members }),
        {
            headers: { 'Content-Type': 'application/json' },
        }
    );
    check(teamRes, { 'team created': (r) => r.status === 201 });

    return { teamName: teamName, userIds: members.map(m => m.user_id) };
}

export default function (data) {
    let resTeam = http.get(`${BASE_URL}/team/get?team_name=${data.teamName}`);
    check(resTeam, {
        'team/get is 200': (r) => r.status === 200,
    });

    let resReview = http.get(`${BASE_URL}/users/getReview?user_id=${data.userIds[0]}`);
    check(resReview, {
        'getReview is 200': (r) => r.status === 200,
    });

    let resStats = http.get(`${BASE_URL}/statistics`);
    check(resStats, {
        'statistics is 200': (r) => r.status === 200,
    });

    let resHealth = http.get(`${BASE_URL}/health`);
    check(resHealth, {
        'health is 200': (r) => r.status === 200,
    });

    sleep(0.1);
}