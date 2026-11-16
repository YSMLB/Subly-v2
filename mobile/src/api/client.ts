import AsyncStorage from '@react-native-async-storage/async-storage';
import { Task, ScheduleSlot, ProposedChange, User, AuthResponse } from '../types';

export const API_BASE_URL = 'http://192.168.1.163:8080';

const TOKEN_KEY = 'subly_auth_token';
const USER_KEY = 'subly_auth_user';

let currentToken: string | null = null;
let currentUser: User | null = null;

export async function setAuthSession(token: string, user: User): Promise<void> {
  currentToken = token;
  currentUser = user;
  await AsyncStorage.setItem(TOKEN_KEY, token);
  await AsyncStorage.setItem(USER_KEY, JSON.stringify(user));
}

export async function clearAuthSession(): Promise<void> {
  currentToken = null;
  currentUser = null;
  await AsyncStorage.removeItem(TOKEN_KEY);
  await AsyncStorage.removeItem(USER_KEY);
}

export async function getAuthSession(): Promise<{ token: string | null; user: User | null }> {
  if (!currentToken) {
    currentToken = await AsyncStorage.getItem(TOKEN_KEY);
  }
  if (!currentUser) {
    const userStr = await AsyncStorage.getItem(USER_KEY);
    if (userStr) {
      try {
        currentUser = JSON.parse(userStr);
      } catch {}
    }
  }
  return { token: currentToken, user: currentUser };
}

export async function getAuthHeaders(): Promise<Record<string, string>> {
  const { token } = await getAuthSession();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  return headers;
}

export async function sendVerificationCode(email: string): Promise<{ status: string; message: string }> {
  const res = await fetch(`${API_BASE_URL}/api/auth/send-code`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || 'Ошибка отправки кода');
  }
  return await res.json();
}

export async function registerUser(req: {
  email: string;
  code: string;
  username: string;
  group_name: string;
  password: string;
}): Promise<AuthResponse> {
  const res = await fetch(`${API_BASE_URL}/api/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || 'Ошибка регистрации');
  }
  const data: AuthResponse = await res.json();
  await setAuthSession(data.token, data.user);
  return data;
}

export async function loginUser(email: string, password: string): Promise<AuthResponse> {
  const res = await fetch(`${API_BASE_URL}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || 'Неверный email или пароль');
  }
  const data: AuthResponse = await res.json();
  await setAuthSession(data.token, data.user);
  return data;
}

const todayStr = '2026-11-15';

let tasksCache: Task[] = [
  {
    id: 1,
    user_id: 1,
    title: 'Learn Go',
    status: 'pending',
    is_pinned: true,
    deadline: '14:00',
    category: '💻 Код',
    date: todayStr,
  },
  {
    id: 2,
    user_id: 1,
    title: 'Лабораторная по БД №2',
    status: 'pending',
    deadline: '18:00',
    category: '🎓 Колледж',
    date: todayStr,
  },
  {
    id: 3,
    user_id: 1,
    title: 'Подготовить конспект по Архитектуре ЭВМ',
    status: 'pending',
    deadline: 'Tomorrow',
    category: '🎓 Колледж',
    date: '2026-11-16',
  },
  {
    id: 4,
    user_id: 1,
    title: 'Work on my programming project',
    status: 'completed',
    deadline: '11:30',
    category: '💻 Код',
    date: todayStr,
    completed_at: todayStr,
  },
  {
    id: 5,
    user_id: 1,
    title: 'To do the 2 tasks for LeetCode',
    status: 'completed',
    deadline: '12:00',
    category: '💻 Код',
    date: todayStr,
    completed_at: todayStr,
  },
  {
    id: 6,
    user_id: 1,
    title: 'Подготовка к коллоквиуму по архитектуре ЭВМ',
    status: 'completed',
    deadline: '16:00',
    category: '🎓 Колледж',
    date: '2026-11-14',
    completed_at: '2026-11-14',
  },
  {
    id: 7,
    user_id: 1,
    title: 'Купить тетради и канцелярию',
    status: 'completed',
    deadline: '18:30',
    category: '📌 Личное',
    date: '2026-11-14',
    completed_at: '2026-11-14',
  },
  {
    id: 8,
    user_id: 1,
    title: 'Сдать лабораторную работу по БД №1',
    status: 'completed',
    deadline: '13:00',
    category: '🎓 Колледж',
    date: '2026-11-13',
    completed_at: '2026-11-13',
  },
  {
    id: 9,
    user_id: 1,
    title: 'Практика по ОС Linux в терминале',
    status: 'completed',
    deadline: '15:20',
    category: '💻 Код',
    date: '2026-11-13',
    completed_at: '2026-11-13',
  },
  {
    id: 10,
    user_id: 1,
    title: 'Go to college',
    status: 'completed',
    deadline: '08:30',
    category: '🎓 Колледж',
    date: '2026-11-11',
    completed_at: '2026-11-11',
  },
  {
    id: 11,
    user_id: 1,
    title: 'Read book (Clean Code)',
    status: 'completed',
    deadline: '21:00',
    category: '📌 Личное',
    date: '2026-11-11',
    completed_at: '2026-11-11',
  },
];

const fullWeekSchedule: Record<number, ScheduleSlot[]> = {
  1: [
    {
      id: 201,
      day_of_week: 1,
      order_num: 1,
      subject: 'Архитектура ЭВМ',
      teacher: 'Смирнов В.П.',
      room: 'ауд. 215',
      start_time: '08:30',
      end_time: '10:00',
    },
    {
      id: 202,
      day_of_week: 1,
      order_num: 2,
      subject: 'Высшая математика',
      teacher: 'Иванов А.А.',
      room: 'ауд. 402',
      start_time: '10:15',
      end_time: '11:45',
    },
    {
      id: 203,
      day_of_week: 1,
      order_num: 3,
      subject: 'Программирование на Go',
      teacher: 'Николаев А.Б.',
      room: 'ауд. 308',
      start_time: '12:00',
      end_time: '13:30',
    },
  ],
  2: [
    {
      id: 204,
      day_of_week: 2,
      order_num: 1,
      subject: 'Физика',
      teacher: 'Волкова И.Г.',
      room: 'ауд. 102',
      start_time: '08:30',
      end_time: '10:00',
    },
    {
      id: 205,
      day_of_week: 2,
      order_num: 2,
      subject: 'Базы данных',
      teacher: 'Кузнецов Д.М.',
      room: 'ауд. 37',
      start_time: '10:15',
      end_time: '11:45',
    },
    {
      id: 206,
      day_of_week: 2,
      order_num: 3,
      subject: 'Английский язык',
      teacher: 'Попова Е.С.',
      room: 'ауд. 108',
      start_time: '12:00',
      end_time: '13:30',
    },
  ],
  3: [
    {
      id: 101,
      day_of_week: 3,
      order_num: 1,
      subject: 'Высшая математика',
      teacher: 'Иванов А.А.',
      room: 'ауд. 402',
      start_time: '08:30',
      end_time: '10:00',
      status: 'past',
    },
    {
      id: 102,
      day_of_week: 3,
      order_num: 2,
      subject: 'Архитектура ЭВМ',
      teacher: 'Смирнов В.П.',
      room: 'ауд. 215',
      start_time: '10:15',
      end_time: '11:45',
      status: 'active',
    },
    {
      id: 103,
      day_of_week: 3,
      order_num: 3,
      subject: 'Базы данных',
      teacher: 'Кузнецов Д.М.',
      room: 'ауд. 37 (корпус Б)',
      start_time: '12:00',
      end_time: '13:30',
      is_replacement: true,
      status: 'upcoming',
    },
    {
      id: 104,
      day_of_week: 3,
      order_num: 4,
      subject: 'Английский язык',
      teacher: 'Попова Е.С.',
      room: 'ауд. 108',
      start_time: '13:45',
      end_time: '15:15',
      status: 'upcoming',
    },
  ],
  4: [
    {
      id: 207,
      day_of_week: 4,
      order_num: 1,
      subject: 'Операционные системы',
      teacher: 'Орлов К.С.',
      room: 'ауд. 304',
      start_time: '08:30',
      end_time: '10:00',
    },
    {
      id: 208,
      day_of_week: 4,
      order_num: 2,
      subject: 'Дискретная математика',
      teacher: 'Иванов А.А.',
      room: 'ауд. 405',
      start_time: '10:15',
      end_time: '11:45',
    },
    {
      id: 209,
      day_of_week: 4,
      order_num: 3,
      subject: 'Физкультура',
      teacher: 'Сидоров Н.А.',
      room: 'Спортзал',
      start_time: '12:00',
      end_time: '13:30',
    },
  ],
  5: [
    {
      id: 210,
      day_of_week: 5,
      order_num: 1,
      subject: 'Программирование на Go',
      teacher: 'Николаев А.Б.',
      room: 'ауд. 308',
      start_time: '10:15',
      end_time: '11:45',
    },
    {
      id: 211,
      day_of_week: 5,
      order_num: 2,
      subject: 'Компьютерные сети',
      teacher: 'Морозов А.В.',
      room: 'ауд. 212',
      start_time: '12:00',
      end_time: '13:30',
    },
    {
      id: 212,
      day_of_week: 5,
      order_num: 3,
      subject: 'Лабораторный практикум',
      teacher: 'Николаев А.Б.',
      room: 'ауд. 308',
      start_time: '13:45',
      end_time: '15:15',
    },
  ],
  6: [
    {
      id: 213,
      day_of_week: 6,
      order_num: 1,
      subject: 'История России',
      teacher: 'Павлов В.И.',
      room: 'ауд. 110',
      start_time: '09:00',
      end_time: '10:30',
    },
    {
      id: 214,
      day_of_week: 6,
      order_num: 2,
      subject: 'Философия',
      teacher: 'Семенова О.Н.',
      room: 'ауд. 110',
      start_time: '10:45',
      end_time: '12:15',
    },
  ],
};

const initialMockChanges: ProposedChange[] = [
  {
    id: 6,
    user_id: 1,
    group_name: 'ИТ-31',
    raw_text: 'Замены на 11.11:\nГруппа: ИТ-31\n3 пара: Базы данных, ауд. 37 (корпус Б)',
    proposed_action: 'add_slot',
    status: 'pending',
    created_at: '2026-11-11 10:00',
  },
];

export async function fetchTasks(userId: number = 1): Promise<Task[]> {
  try {
    const headers = await getAuthHeaders();
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 2000);
    const response = await fetch(`${API_BASE_URL}/tasks?user_id=${userId}`, {
      headers,
      signal: controller.signal,
    });
    clearTimeout(timeout);
    if (!response.ok) throw new Error('API error');
    const data = await response.json();
    if (Array.isArray(data) && data.length > 0) {
      tasksCache = data;
    }
  } catch {
    // Offline / fallback cache
  }
  return [...tasksCache];
}

export async function updateTaskStatus(taskId: number, status: 'pending' | 'completed'): Promise<void> {
  const currentToday = new Date().toISOString().split('T')[0];
  tasksCache = tasksCache.map((t) =>
    t.id === taskId
      ? {
          ...t,
          status,
          completed_at: status === 'completed' ? currentToday : undefined,
        }
      : t
  );

  try {
    const headers = await getAuthHeaders();
    await fetch(`${API_BASE_URL}/tasks/${taskId}/status`, {
      method: 'PATCH',
      headers,
      body: JSON.stringify({ status }),
    });
  } catch (e) {
    console.log('Status update offline fallback');
  }
}

export interface CreateTaskOptions {
  deadline?: string;
  category?: string;
  is_pinned?: boolean;
  date?: string;
}

export async function createTask(
  userId: number,
  title: string,
  options?: CreateTaskOptions
): Promise<Task> {
  const currentToday = new Date().toISOString().split('T')[0];
  const newTask: Task = {
    id: Date.now(),
    user_id: userId,
    title,
    status: 'pending',
    deadline: options?.deadline || 'Today',
    category: options?.category,
    is_pinned: options?.is_pinned || false,
    date: options?.date || currentToday,
  };

  tasksCache = [newTask, ...tasksCache];

  try {
    const headers = await getAuthHeaders();
    await fetch(`${API_BASE_URL}/tasks`, {
      method: 'POST',
      headers,
      body: JSON.stringify({
        user_id: userId,
        title,
        status: 'pending',
      }),
    });
  } catch (e) {
    console.log('Create task offline fallback');
  }

  return newTask;
}

export async function deleteTask(taskId: number): Promise<void> {
  tasksCache = tasksCache.filter((t) => t.id !== taskId);
  try {
    const headers = await getAuthHeaders();
    await fetch(`${API_BASE_URL}/tasks/${taskId}`, {
      method: 'DELETE',
      headers,
    });
  } catch (e) {
    console.log('Delete task offline fallback');
  }
}

export async function fetchSchedule(dayOfWeek: number = 3): Promise<ScheduleSlot[]> {
  try {
    const headers = await getAuthHeaders();
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 2000);
    const response = await fetch(`${API_BASE_URL}/schedule/today?day_of_week=${dayOfWeek}`, {
      headers,
      signal: controller.signal,
    });
    clearTimeout(timeout);
    if (!response.ok) throw new Error('API error');
    const data = await response.json();
    return Array.isArray(data) && data.length > 0 ? data : fullWeekSchedule[dayOfWeek] || [];
  } catch {
    return fullWeekSchedule[dayOfWeek] || [];
  }
}

export async function fetchPendingChanges(userId: number = 1): Promise<ProposedChange[]> {
  try {
    const headers = await getAuthHeaders();
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 2000);
    const response = await fetch(`${API_BASE_URL}/changes/pending?user_id=${userId}`, {
      headers,
      signal: controller.signal,
    });
    clearTimeout(timeout);
    if (!response.ok) throw new Error('API error');
    const data = await response.json();
    return Array.isArray(data) && data.length > 0 ? data : initialMockChanges;
  } catch {
    return initialMockChanges;
  }
}

export async function reviewChange(id: number, status: 'approved' | 'rejected'): Promise<void> {
  try {
    const headers = await getAuthHeaders();
    await fetch(`${API_BASE_URL}/changes/${id}/status`, {
      method: 'PATCH',
      headers,
      body: JSON.stringify({ status, reviewed_by: 1 }),
    });
  } catch (e) {
    console.log('Review offline fallback');
  }
}

