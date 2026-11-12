import { Task, ScheduleSlot, ProposedChange } from '../types';

export const API_BASE_URL = 'http://192.168.0.103:8080';

const initialMockTasks: Task[] = [
  {
    id: 1,
    user_id: 1,
    title: 'Learn Go',
    status: 'pending',
    is_pinned: true,
    deadline: 'Today',
  },
  {
    id: 2,
    user_id: 1,
    title: 'Лабораторная по БД №2',
    status: 'pending',
    deadline: 'Today',
  },
  {
    id: 3,
    user_id: 1,
    title: 'Подготовить конспект по Архитектуре ЭВМ',
    status: 'pending',
    deadline: 'Tomorrow',
  },
  {
    id: 4,
    user_id: 1,
    title: 'Work on my programming project',
    status: 'completed',
    deadline: 'Today',
  },
  {
    id: 5,
    user_id: 1,
    title: 'To do the 2 tasks for LeetCode',
    status: 'completed',
    deadline: 'Today',
  },
];

const initialMockSchedule: ScheduleSlot[] = [
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
];

const initialMockChanges: ProposedChange[] = [
  {
    id: 6,
    user_id: 1,
    group_name: 'ИТ-31',
    raw_text: 'Замены на 11.11:\nГруппа: ИТ-31\n3 пара: Базы данных, ауд. 37',
    proposed_action: 'add_slot',
    status: 'pending',
    created_at: '2026-11-11 10:00',
  },
];

export async function fetchTasks(userId: number = 1): Promise<Task[]> {
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 2000);
    const response = await fetch(`${API_BASE_URL}/tasks?user_id=${userId}`, {
      signal: controller.signal,
    });
    clearTimeout(timeout);
    if (!response.ok) throw new Error('API error');
    const data = await response.json();
    return Array.isArray(data) && data.length > 0 ? data : initialMockTasks;
  } catch {
    return initialMockTasks;
  }
}

export async function updateTaskStatus(taskId: number, status: 'pending' | 'completed'): Promise<void> {
  try {
    await fetch(`${API_BASE_URL}/tasks/${taskId}/status`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status }),
    });
  } catch (e) {
    console.log('Status update offline fallback');
  }
}

export async function createTask(userId: number, title: string): Promise<Task> {
  const newTask: Task = {
    id: Date.now(),
    user_id: userId,
    title,
    status: 'pending',
    deadline: 'Today',
  };

  try {
    await fetch(`${API_BASE_URL}/tasks`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
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

export async function fetchSchedule(userId: number = 1): Promise<ScheduleSlot[]> {
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 2000);
    const response = await fetch(`${API_BASE_URL}/schedule/today?user_id=${userId}`, {
      signal: controller.signal,
    });
    clearTimeout(timeout);
    if (!response.ok) throw new Error('API error');
    const data = await response.json();
    return Array.isArray(data) && data.length > 0 ? data : initialMockSchedule;
  } catch {
    return initialMockSchedule;
  }
}

export async function fetchPendingChanges(userId: number = 1): Promise<ProposedChange[]> {
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 2000);
    const response = await fetch(`${API_BASE_URL}/changes/pending?user_id=${userId}`, {
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
    await fetch(`${API_BASE_URL}/changes/${id}/status`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status, reviewed_by: 1 }),
    });
  } catch (e) {
    console.log('Review offline fallback');
  }
}
