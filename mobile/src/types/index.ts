export interface Task {
  id: number;
  user_id: number;
  title: string;
  description?: string;
  status: 'pending' | 'in_progress' | 'completed';
  deadline?: string;
  is_pinned?: boolean;
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

export type TabType = 'tasks' | 'schedule' | 'moderation';
