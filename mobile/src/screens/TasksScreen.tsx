import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  Modal,
  TextInput,
  KeyboardAvoidingView,
  Platform,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Task, ScheduleSlot } from '../types';
import { fetchTasks, updateTaskStatus, createTask, deleteTask, fetchSchedule } from '../api/client';
import { sendPushNotification } from '../services/notificationService';

interface TasksScreenProps {
  onNavigateToSchedule?: () => void;
  onTriggerNotification?: (title: string, message: string) => void;
}

export const TasksScreen: React.FC<TasksScreenProps> = ({
  onNavigateToSchedule,
  onTriggerNotification,
}) => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [todaySchedule, setTodaySchedule] = useState<ScheduleSlot[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isScheduleModalOpen, setIsScheduleModalOpen] = useState(false);
  const [newTitle, setNewTitle] = useState('');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    const [taskList, scheduleList] = await Promise.all([
      fetchTasks(),
      fetchSchedule(),
    ]);

    setTodaySchedule(scheduleList);

    // Автоматическая задача на учебный день (Пн - Сб)
    const currentDay = new Date().getDay(); // 0 - Вс, 1 - Пн ... 6 - Сб
    const isSchoolDay = currentDay >= 1 && currentDay <= 6;
    const hasReplacement = scheduleList.some((s) => s.is_replacement);

    let combinedTasks = [...taskList];

    if (isSchoolDay && !combinedTasks.some((t) => t.is_auto_college)) {
      const collegeTask: Task = {
        id: 999999,
        user_id: 1,
        title: '🎓 Пары в колледже (Расписание на сегодня)',
        status: 'pending',
        is_pinned: true,
        is_auto_college: true,
        deadline: hasReplacement ? '⚡ Есть замена' : '4 пары',
      };
      combinedTasks = [collegeTask, ...combinedTasks];
    }

    setTasks(combinedTasks);
  };

  const toggleTask = (task: Task) => {
    const newStatus = task.status === 'completed' ? 'pending' : 'completed';
    setTasks((prev) =>
      prev.map((t) => (t.id === task.id ? { ...t, status: newStatus } : t))
    );
    if (!task.is_auto_college) {
      updateTaskStatus(task.id, newStatus);
    }
  };

  const togglePin = (task: Task) => {
    setTasks((prev) =>
      prev.map((t) => (t.id === task.id ? { ...t, is_pinned: !t.is_pinned } : t))
    );
  };

  const handleDeleteTask = (task: Task) => {
    Alert.alert('Удалить задачу?', `«${task.title}»`, [
      { text: 'Отмена', style: 'cancel' },
      {
        text: 'Удалить',
        style: 'destructive',
        onPress: async () => {
          setTasks((prev) => prev.filter((t) => t.id !== task.id));
          if (!task.is_auto_college) {
            await deleteTask(task.id);
          }
        },
      },
    ]);
  };

  const handleCreateTask = async () => {
    if (!newTitle.trim()) return;
    const created = await createTask(1, newTitle.trim());
    setTasks((prev) => [created, ...prev]);
    setNewTitle('');
    setIsModalOpen(false);

    if (onTriggerNotification) {
      onTriggerNotification('Новая задача создана', `«${created.title}» добавлена в список дел`);
    } else {
      sendPushNotification('Новая задача создана', `«${created.title}» добавлена в список дел`);
    }
  };

  const pinnedTasks = tasks.filter((t) => t.is_pinned && t.status !== 'completed');
  const todayTasks = tasks.filter((t) => !t.is_pinned && t.status !== 'completed');
  const completedTasks = tasks.filter((t) => t.status === 'completed');

  const renderTaskItem = (task: Task, isPinnedView: boolean = false) => {
    const isCollege = task.is_auto_college;

    return (
      <View
        key={task.id}
        style={[
          styles.taskRow,
          isPinnedView && styles.pinnedRowBorder,
        ]}
      >
        <TouchableOpacity
          style={styles.checkbox}
          onPress={() => toggleTask(task)}
        >
          <Ionicons
            name={task.status === 'completed' ? 'checkbox' : 'square-outline'}
            size={22}
            color={task.status === 'completed' ? '#00E5FF' : '#7E8696'}
          />
        </TouchableOpacity>

        <TouchableOpacity
          style={styles.taskContent}
          activeOpacity={isCollege ? 0.7 : 1}
          onPress={() => {
            if (isCollege) {
              setIsScheduleModalOpen(true);
            }
          }}
        >
          <Text
            style={[
              styles.taskTitle,
              task.status === 'completed' && styles.completedTitle,
              isCollege && styles.collegeTaskTitle,
            ]}
          >
            {task.title}
          </Text>

          {task.deadline && (
            <View style={[styles.badge, isCollege && styles.collegeBadge]}>
              <Text style={[styles.badgeText, isCollege && styles.collegeBadgeText]}>
                {task.deadline}
              </Text>
            </View>
          )}
        </TouchableOpacity>

        <View style={styles.rowActions}>
          <TouchableOpacity
            style={styles.actionIcon}
            onPress={() => togglePin(task)}
            hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
          >
            <Ionicons
              name={task.is_pinned ? 'pin' : 'pin-outline'}
              size={18}
              color={task.is_pinned ? '#00E5FF' : '#5E6778'}
            />
          </TouchableOpacity>

          {!isCollege && (
            <TouchableOpacity
              style={styles.actionIcon}
              onPress={() => handleDeleteTask(task)}
              hitSlop={{ top: 8, bottom: 8, left: 8, right: 8 }}
            >
              <Ionicons name="trash-outline" size={18} color="#7E8696" />
            </TouchableOpacity>
          )}
        </View>
      </View>
    );
  };

  return (
    <View style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        <View style={styles.header}>
          <View>
            <Text style={styles.headerCategory}>Subly / Задачи</Text>
            <Text style={styles.headerTitle}>Today</Text>
          </View>
          <TouchableOpacity
            style={styles.testNotifyBtn}
            onPress={() => {
              const msg = '3 пара: Базы данных (ауд. 37, корпус Б) вместо Физики';
              if (onTriggerNotification) {
                onTriggerNotification('⚡ Замена в расписании!', msg);
              } else {
                sendPushNotification('⚡ Замена в расписании!', msg);
              }
            }}
          >
            <Ionicons name="notifications-outline" size={20} color="#00E5FF" />
          </TouchableOpacity>
        </View>

        {pinnedTasks.length > 0 && (
          <View style={styles.section}>
            <View style={styles.sectionHeader}>
              <Ionicons name="pin" size={16} color="#00E5FF" />
              <Text style={styles.sectionTitle}>Pinned</Text>
              <Text style={styles.pinnedCount}>({pinnedTasks.length})</Text>
            </View>
            <View style={styles.pinnedCardsContainer}>
              {pinnedTasks.map((task) => renderTaskItem(task, true))}
            </View>
          </View>
        )}

        <View style={styles.section}>
          <View style={styles.sectionHeader}>
            <Ionicons name="today-outline" size={16} color="#7E8696" />
            <Text style={styles.sectionTitle}>Today</Text>
            <Text style={styles.sectionCount}>{todayTasks.length}</Text>
          </View>
          <View style={styles.card}>
            {todayTasks.length === 0 ? (
              <Text style={styles.emptyText}>Все задачи выполнены 🎉</Text>
            ) : (
              todayTasks.map((task) => renderTaskItem(task, false))
            )}
          </View>
        </View>

        {completedTasks.length > 0 && (
          <View style={styles.section}>
            <View style={styles.sectionHeader}>
              <Ionicons name="checkmark-done" size={16} color="#00E5FF" />
              <Text style={styles.sectionTitle}>Completed</Text>
              <Text style={styles.sectionCount}>{completedTasks.length}</Text>
            </View>
            <View style={styles.card}>
              {completedTasks.map((task) => renderTaskItem(task, false))}
            </View>
          </View>
        )}
      </ScrollView>

      <TouchableOpacity
        style={styles.fab}
        activeOpacity={0.8}
        onPress={() => setIsModalOpen(true)}
      >
        <Ionicons name="add" size={32} color="#000000" />
      </TouchableOpacity>

      {/* Модалка создания новой задачи */}
      <Modal
        visible={isModalOpen}
        transparent
        animationType="fade"
        onRequestClose={() => setIsModalOpen(false)}
      >
        <KeyboardAvoidingView
          behavior={Platform.OS === 'ios' ? 'padding' : undefined}
          style={styles.modalOverlay}
        >
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>Новая задача</Text>
            <TextInput
              style={styles.modalInput}
              placeholder="Что нужно сделать?"
              placeholderTextColor="#717A8C"
              value={newTitle}
              onChangeText={setNewTitle}
              autoFocus
            />
            <View style={styles.modalButtons}>
              <TouchableOpacity
                style={styles.modalCancel}
                onPress={() => setIsModalOpen(false)}
              >
                <Text style={styles.modalCancelText}>Отмена</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={styles.modalSubmit}
                onPress={handleCreateTask}
              >
                <Text style={styles.modalSubmitText}>Добавить</Text>
              </TouchableOpacity>
            </View>
          </View>
        </KeyboardAvoidingView>
      </Modal>

      {/* Шторка просмотра расписания из авто-задачи */}
      <Modal
        visible={isScheduleModalOpen}
        transparent
        animationType="slide"
        onRequestClose={() => setIsScheduleModalOpen(false)}
      >
        <View style={styles.scheduleModalOverlay}>
          <View style={styles.scheduleModalSheet}>
            <View style={styles.sheetHandle} />
            <View style={styles.sheetHeader}>
              <View>
                <Text style={styles.sheetTitle}>Пары на сегодня</Text>
                <Text style={styles.sheetSubtitle}>Итоговое расписание с заменами</Text>
              </View>
              <TouchableOpacity
                onPress={() => setIsScheduleModalOpen(false)}
                hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
              >
                <Ionicons name="close-circle" size={26} color="#7E8696" />
              </TouchableOpacity>
            </View>

            <ScrollView showsVerticalScrollIndicator={false} style={styles.sheetScroll}>
              {todaySchedule.map((slot) => {
                const isRep = slot.is_replacement;
                return (
                  <View
                    key={slot.id}
                    style={[styles.sheetSlot, isRep && styles.sheetSlotRep]}
                  >
                    <View style={styles.sheetSlotTop}>
                      <Text style={[styles.sheetOrder, isRep && styles.sheetOrderRep]}>
                        {slot.order_num} пара • {slot.start_time} - {slot.end_time}
                      </Text>
                      {isRep && (
                        <View style={styles.repPill}>
                          <Text style={styles.repPillText}>Замена</Text>
                        </View>
                      )}
                    </View>
                    <Text style={styles.sheetSubject}>{slot.subject}</Text>
                    <View style={styles.sheetSlotBottom}>
                      <Text style={styles.sheetMeta}>📍 {slot.room}</Text>
                      {slot.teacher && <Text style={styles.sheetMeta}>👤 {slot.teacher}</Text>}
                    </View>
                  </View>
                );
              })}
            </ScrollView>

            <TouchableOpacity
              style={styles.sheetFullBtn}
              onPress={() => {
                setIsScheduleModalOpen(false);
                if (onNavigateToSchedule) onNavigateToSchedule();
              }}
            >
              <Text style={styles.sheetFullBtnText}>Перейти во вкладку «Пары»</Text>
              <Ionicons name="arrow-forward" size={18} color="#000000" />
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000000',
  },
  scrollContent: {
    paddingHorizontal: 20,
    paddingTop: 16,
    paddingBottom: 110,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 20,
  },
  headerCategory: {
    color: '#7E8696',
    fontSize: 13,
    fontWeight: '600',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  headerTitle: {
    color: '#FFFFFF',
    fontSize: 38,
    fontWeight: '800',
    marginTop: 4,
    letterSpacing: -0.5,
  },
  testNotifyBtn: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: '#121622',
    borderWidth: 1,
    borderColor: '#1E2536',
    justifyContent: 'center',
    alignItems: 'center',
  },
  section: {
    marginBottom: 22,
  },
  sectionHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 10,
    paddingHorizontal: 4,
  },
  sectionTitle: {
    color: '#CAD1DE',
    fontSize: 15,
    fontWeight: '700',
    marginLeft: 8,
  },
  sectionCount: {
    color: '#7E8696',
    fontSize: 13,
    marginLeft: 8,
    fontWeight: '600',
  },
  pinnedCount: {
    color: '#00E5FF',
    fontSize: 13,
    marginLeft: 6,
    fontWeight: '700',
  },
  pinnedCardsContainer: {
    gap: 10,
  },
  card: {
    backgroundColor: '#121622',
    borderRadius: 22,
    paddingHorizontal: 16,
    paddingVertical: 4,
    borderWidth: 1,
    borderColor: '#1E2536',
  },
  taskRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 14,
    paddingHorizontal: 16,
    backgroundColor: '#121622',
    borderRadius: 18,
    borderWidth: 1,
    borderColor: '#1E2536',
  },
  pinnedRowBorder: {
    borderColor: '#00E5FF',
    borderWidth: 1.5,
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.65,
    shadowRadius: 10,
    elevation: 6,
  },
  checkbox: {
    marginRight: 12,
  },
  taskContent: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingRight: 8,
  },
  taskTitle: {
    flex: 1,
    color: '#FFFFFF',
    fontSize: 15,
    fontWeight: '500',
  },
  collegeTaskTitle: {
    fontWeight: '700',
    color: '#E0F7FA',
  },
  completedTitle: {
    color: '#636C7E',
    textDecorationLine: 'line-through',
  },
  badge: {
    backgroundColor: 'rgba(0, 229, 255, 0.12)',
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 10,
    borderWidth: 1,
    borderColor: 'rgba(0, 229, 255, 0.3)',
    marginLeft: 8,
  },
  badgeText: {
    color: '#00E5FF',
    fontSize: 11,
    fontWeight: '600',
  },
  collegeBadge: {
    backgroundColor: 'rgba(0, 229, 255, 0.2)',
    borderColor: '#00E5FF',
  },
  collegeBadgeText: {
    color: '#00E5FF',
    fontWeight: '700',
  },
  rowActions: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  actionIcon: {
    padding: 4,
  },
  emptyText: {
    color: '#7E8696',
    fontSize: 14,
    textAlign: 'center',
    paddingVertical: 18,
  },
  fab: {
    position: 'absolute',
    right: 24,
    bottom: 100,
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: '#00E5FF',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.7,
    shadowRadius: 14,
    elevation: 8,
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.75)',
    justifyContent: 'center',
    paddingHorizontal: 24,
  },
  modalContent: {
    backgroundColor: '#121622',
    borderRadius: 24,
    padding: 24,
    borderWidth: 1,
    borderColor: '#222B3D',
  },
  modalTitle: {
    color: '#FFFFFF',
    fontSize: 20,
    fontWeight: '700',
    marginBottom: 16,
  },
  modalInput: {
    backgroundColor: '#171D2B',
    borderRadius: 14,
    paddingHorizontal: 16,
    height: 52,
    color: '#FFFFFF',
    fontSize: 16,
    borderWidth: 1,
    borderColor: '#263145',
    marginBottom: 20,
  },
  modalButtons: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 12,
  },
  modalCancel: {
    paddingVertical: 10,
    paddingHorizontal: 18,
    borderRadius: 12,
  },
  modalCancelText: {
    color: '#8A93A6',
    fontSize: 15,
    fontWeight: '600',
  },
  modalSubmit: {
    backgroundColor: '#00E5FF',
    paddingVertical: 10,
    paddingHorizontal: 22,
    borderRadius: 12,
  },
  modalSubmitText: {
    color: '#000000',
    fontSize: 15,
    fontWeight: '700',
  },
  scheduleModalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.8)',
    justifyContent: 'flex-end',
  },
  scheduleModalSheet: {
    backgroundColor: '#121622',
    borderTopLeftRadius: 32,
    borderTopRightRadius: 32,
    paddingHorizontal: 20,
    paddingTop: 12,
    paddingBottom: 36,
    maxHeight: '75%',
    borderWidth: 1,
    borderColor: '#222B3D',
  },
  sheetHandle: {
    width: 40,
    height: 4,
    backgroundColor: '#353D50',
    borderRadius: 2,
    alignSelf: 'center',
    marginBottom: 16,
  },
  sheetHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 16,
  },
  sheetTitle: {
    color: '#FFFFFF',
    fontSize: 22,
    fontWeight: '800',
  },
  sheetSubtitle: {
    color: '#8A93A6',
    fontSize: 13,
    marginTop: 2,
  },
  sheetScroll: {
    marginBottom: 16,
  },
  sheetSlot: {
    backgroundColor: '#171D2B',
    borderRadius: 16,
    padding: 14,
    marginBottom: 10,
    borderWidth: 1,
    borderColor: '#222B3D',
  },
  sheetSlotRep: {
    borderColor: '#00E5FF',
    backgroundColor: '#111A26',
  },
  sheetSlotTop: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 6,
  },
  sheetOrder: {
    color: '#8A93A6',
    fontSize: 12,
    fontWeight: '600',
  },
  sheetOrderRep: {
    color: '#00E5FF',
  },
  repPill: {
    backgroundColor: 'rgba(0, 229, 255, 0.15)',
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 6,
  },
  repPillText: {
    color: '#00E5FF',
    fontSize: 10,
    fontWeight: '700',
  },
  sheetSubject: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '700',
    marginBottom: 6,
  },
  sheetSlotBottom: {
    flexDirection: 'row',
    gap: 14,
  },
  sheetMeta: {
    color: '#8A93A6',
    fontSize: 12,
  },
  sheetFullBtn: {
    flexDirection: 'row',
    backgroundColor: '#00E5FF',
    height: 52,
    borderRadius: 16,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 8,
  },
  sheetFullBtnText: {
    color: '#000000',
    fontSize: 15,
    fontWeight: '700',
  },
});
