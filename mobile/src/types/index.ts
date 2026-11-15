export interface Task {
  id: number;
  user_id: number;
  title: string;
  description?: string;
  status: 'pending' | 'in_progress' | 'completed';
  deadline?: string;
  category?: string;
  is_pinned?: boolean;
  is_auto_college?: boolean;
  date?: string; // YYYY-MM-DD
  completed_at?: string; // YYYY-MM-DD
  created_at?: string;
}

export interface ScheduleSlot {
  id: number;
  day_of_week: number;
  order_num: number;
  subject: string;
  teacher?: string;
  room: string;
  start_time: string;
  end_time: string;
  is_replacement?: boolean;
  status?: 'past' | 'active' | 'upcoming';
}

export interface ProposedChange {
  id: number;
  user_id: number;
  group_name: string;
  raw_text: string;
  proposed_action: string;
  status: 'pending' | 'approved' | 'rejected';
  created_at?: string;
}

export interface NotificationPayload {
  id: string;
  title: string;
  message: string;
  type: 'replacement' | 'reminder' | 'system' | 'task';
  timestamp?: string;
}

export type TabType = 'tasks' | 'schedule' | 'moderation';
