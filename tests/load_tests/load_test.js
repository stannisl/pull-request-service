import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Кастомные метрики
const teamCreationDuration = new Trend('team_creation_duration');
const teamRetrievalDuration = new Trend('team_retrieval_duration');
const errorRate = new Rate('errors');

export const options = {
    stages: [
        // Плавный рост нагрузки до 5 RPS (ваш целевой показатель)
        { duration: '1m', target: 5 },    // 5 пользователей = ~5 RPS
        { duration: '3m', target: 5 },    // Стабильная нагрузка 5 RPS
        { duration: '1m', target: 10 },   // Пиковая нагрузка
        { duration: '1m', target: 5 },    // Возврат к нормальной нагрузке
        { duration: '1m', target: 0 },    // Завершение
    ],
    thresholds: {
        // Основные SLI
        http_req_duration: ['p(95)<300'],  // 95% запросов < 300ms
        http_req_failed: ['rate<0.001'],   // < 0.1% ошибок
        errors: ['rate<0.001'],            // < 0.1% кастомных ошибок

        // Дополнительные метрики для мониторинга
        team_creation_duration: ['p(95)<250'],
        team_retrieval_duration: ['p(95)<200'],
    },
};

const BASE_URL = 'http://localhost:8080';

// Глобальная инициализация (выполняется один раз)
export function setup() {
    console.log('Starting load test with target: 5 RPS, 300ms p95 latency');
    return { startTime: new Date().toISOString() };
}

export default function (data) {
    const timestamp = Date.now();
    const randomId = Math.floor(Math.random() * 10000);
    const teamName = `load_test_team_${timestamp}_${__VU}_${randomId}`;

    // Сценарий 1: Создание команды
    const teamData = JSON.stringify({
        team_name: teamName,
        members: [
            {
                user_id: `user1_${timestamp}_${__VU}_${randomId}`,
                username: `Alice_${timestamp}`,
                is_active: true
            },
            {
                user_id: `user2_${timestamp}_${__VU}_${randomId}`,
                username: `Bob_${timestamp}`,
                is_active: true
            },
            {
                user_id: `user3_${timestamp}_${__VU}_${randomId}`,
                username: `Charlie_${timestamp}`,
                is_active: true
            }
        ]
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'User-Agent': 'k6-load-test/1.0'
        },
        tags: {
            name: 'create_team',
            endpoint: '/team/add',
            method: 'POST'
        },
    };

    // Запрос на создание команды
    let createStart = Date.now();
    let createRes = http.post(`${BASE_URL}/team/add`, teamData, params);
    let createDuration = Date.now() - createStart;
    teamCreationDuration.add(createDuration);

    let createCheck = check(createRes, {
        'team created successfully': (r) => r.status === 201,
        'response time acceptable': (r) => r.timings.duration < 300,
        'has correct content type': (r) => r.headers['Content-Type']?.includes('application/json'),
    });

    if (!createCheck) {
        errorRate.add(1);
        console.error(`Team creation failed: ${createRes.status} - ${createRes.body}`);
    }

    // Короткая пауза между запросами
    sleep(0.1);

    // Сценарий 2: Получение команды
    let getStart = Date.now();
    let getRes = http.get(`${BASE_URL}/team/get?team_name=${teamName}`, {
        tags: {
            name: 'get_team',
            endpoint: '/team/get',
            method: 'GET'
        }
    });
    let getDuration = Date.now() - getStart;
    teamRetrievalDuration.add(getDuration);

    let getCheck = check(getRes, {
        'team retrieved successfully': (r) => r.status === 200,
        'team data is correct': (r) => {
            try {
                const data = JSON.parse(r.body);
                return data.team && data.team.team_name === teamName;
            } catch {
                return false;
            }
        },
        'response time acceptable': (r) => r.timings.duration < 300,
    });

    if (!getCheck) {
        errorRate.add(1);
        console.error(`Team retrieval failed: ${getRes.status} - ${getRes.body}`);
    }

    // Имитация реального пользовательского поведения
    sleep(Math.random() * 0.3 + 0.1); // Случайная пауза 100-400ms
}

export function teardown(data) {
    console.log(`Load test completed. Started at: ${data.startTime}`);
}