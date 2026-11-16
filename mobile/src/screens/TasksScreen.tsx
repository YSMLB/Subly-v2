import React, { useState, useEffect, useRef } from 'react';
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
  Linking,
  Animated,
  Dimensions,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as Haptics from 'expo-haptics';
import { Task, ScheduleSlot, User } from '../types';
import { fetchTasks, updateTaskStatus, createTask, deleteTask, fetchSchedule, getAuthSession } from '../api/client';
import { sendPushNotification } from '../services/notificationService';

interface TasksScreenProps {
  onNavigateToSchedule?: () => void;
  onTriggerNotification?: (title: string, message: string, type?: 'replacement' | 'reminder' | 'system' | 'task') => void;
  onLogout?: () => void;
}

const CATEGORIES = ['🎓 Колледж', '💻 Код', '⚡ Срочно', '📌 Личное'];
const QUICK_TIMES = ['10:00', '14:00', '18:00', '21:00', 'Весь день'];

const { width: SCREEN_WIDTH } = Dimensions.get('window');
const DRAWER_WIDTH = Math.min(SCREEN_WIDTH * 0.78, 300);

export const TasksScreen: React.FC<TasksScreenProps> = ({
  onNavigateToSchedule,
  onTriggerNotification,
  onLogout,
}) => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [todaySchedule, setTodaySchedule] = useState<ScheduleSlot[]>([]);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isScheduleModalOpen, setIsScheduleModalOpen] = useState(false);
  const [isBurgerMenuOpen, setIsBurgerMenuOpen] = useState(false);

  // Анимация бокового меню слева
  const slideAnim = useRef(new Animated.Value(-DRAWER_WIDTH)).current;

  const openDrawer = () => {
    setIsBurgerMenuOpen(true);
    Animated.spring(slideAnim, {
      toValue: 0,
      useNativeDriver: true,
      bounciness: 4,
      speed: 14,
    }).start();
  };

  const closeDrawer = (callback?: () => void) => {
    Animated.timing(slideAnim, {
      toValue: -DRAWER_WIDTH,
      duration: 200,
      useNativeDriver: true,
    }).start(() => {
      setIsBurgerMenuOpen(false);
      if (callback) callback();
    });
  };

  // Поля модалки создания задачи
  const [newTitle, setNewTitle] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('🎓 Колледж');
  const [selectedTime, setSelectedTime] = useState<string>('14:00');
  const [customTime, setCustomTime] = useState<string>('');
  const [isNewPinned, setIsNewPinned] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    const [{ user }, taskList, scheduleList] = await Promise.all([
      getAuthSession(),
      fetchTasks(),
      fetchSchedule(new Date().getDay() || 1),
    ]);

    if (user) {
      setCurrentUser(user);
    }

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
        category: '🎓 Колледж',
        deadline: hasReplacement ? '⚡ Есть замена' : 'Расписание',
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

    const timeLabel = customTime.trim() || selectedTime || 'Today';

    const created = await createTask(1, newTitle.trim(), {
      deadline: timeLabel,
      category: selectedCategory,
      is_pinned: isNewPinned,
      date: new Date().toISOString().split('T')[0],
    });

    setTasks((prev) => [created, ...prev]);

    // Сброс формы
    setNewTitle('');
    setCustomTime('');
    setSelectedTime('14:00');
    setIsNewPinned(false);
    setIsModalOpen(false);

    if (onTriggerNotification) {
      onTriggerNotification(
        'Новая задача создана',
        `«${created.title}» добавлена в список дел`,
        'task'
      );
    } else {
      sendPushNotification('Subly', 'Есть изменения', true);
    }
  };

  const handleContactDev = async () => {
    await Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    Linking.openURL('https://t.me/amiri2sb').catch(() => {
      Alert.alert(
        'Связь с разработчиком',
        'Telegram разработчика: @amiri2sb\nПочта поддержки: support@subly.app',
        [{ text: 'Понятно' }]
      );
    });
  };

  const handleLogoutPress = () => {
    closeDrawer(() => {
      Alert.alert(
        'Выйти из аккаунта?',
        'Вы сможете войти снова в любое время.',
        [
          { text: 'Отмена', style: 'cancel' },
          {
            text: 'Выйти',
            style: 'destructive',
            onPress: () => {
              if (onLogout) onLogout();
            },
          },
        ]
      );
    });
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
          <View style={styles.titleWrapper}>
            <Text
              style={[
                styles.taskTitle,
                task.status === 'completed' && styles.completedTitle,
                isCollege && styles.collegeTaskTitle,
              ]}
            >
              {task.title}
            </Text>

            <View style={styles.badgesRow}>
              {task.category && (
                <View style={styles.categoryBadge}>
                  <Text style={styles.categoryBadgeText}>{task.category}</Text>
                </View>
              )}
              {task.deadline && (
                <View style={[styles.badge, isCollege && styles.collegeBadge]}>
                  <Ionicons name="time-outline" size={11} color="#00E5FF" style={{ marginRight: 3 }} />
                  <Text style={[styles.badgeText, isCollege && styles.collegeBadgeText]}>
                    {task.deadline}
                  </Text>
                </View>
              )}
            </View>
          </View>
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
          <View style={styles.headerLeft}>
            <TouchableOpacity
              style={styles.burgerBtn}
              activeOpacity={0.7}
              onPress={openDrawer}
            >
              <Ionicons name="menu" size={24} color="#FFFFFF" />
            </TouchableOpacity>

            <View style={styles.headerTitleBox}>
              <Text style={styles.headerCategory}>Subly / Задачи</Text>
              <Text style={styles.headerTitle}>Today</Text>
            </View>
          </View>

          <TouchableOpacity
            style={styles.testNotifyBtn}
            activeOpacity={0.7}
            onPress={() => {
              const msg = '3 пара: Базы данных (ауд. 37, корпус Б) вместо Физики';
              if (onTriggerNotification) {
                onTriggerNotification('⚡ Замена в расписании!', msg, 'replacement');
              } else {
                sendPushNotification('Subly', 'Есть изменения', true);
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
          <View style={styles.tasksListContainer}>
            {todayTasks.length === 0 ? (
              <View style={styles.emptyCard}>
                <Text style={styles.emptyText}>Все задачи выполнены 🎉</Text>
              </View>
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
            <View style={styles.tasksListContainer}>
              {completedTasks.map((task) => renderTaskItem(task, false))}
            </View>
          </View>
        )}
      </ScrollView>

      {/* Floating Action Button */}
      <TouchableOpacity
        style={styles.fab}
        activeOpacity={0.85}
        onPress={() => setIsModalOpen(true)}
      >
        <Ionicons name="add" size={32} color="#000000" />
      </TouchableOpacity>

      {/* =================== БОКОВОЕ БУРГЕР МЕНЮ (СЛЕВА) =================== */}
      <Modal
        visible={isBurgerMenuOpen}
        transparent
        animationType="none"
        onRequestClose={() => closeDrawer()}
      >
        <View style={styles.burgerOverlay}>
          {/* Фон затемнения справа */}
          <TouchableOpacity
            style={styles.burgerBackdrop}
            activeOpacity={1}
            onPress={() => closeDrawer()}
          />

          {/* Панель, выезжающая слева */}
          <Animated.View
            style={[
              styles.drawerContainer,
              {
                width: DRAWER_WIDTH,
                transform: [{ translateX: slideAnim }],
              },
            ]}
          >
            {/* Header профиля без лишней залупы */}
            <View style={styles.drawerHeader}>
              <View style={styles.drawerProfileRow}>
                <View style={styles.drawerAvatar}>
                  <Text style={styles.drawerAvatarText}>
                    {currentUser?.username ? currentUser.username.slice(0, 2).toUpperCase() : 'YS'}
                  </Text>
                </View>
                <View style={styles.drawerUserInfo}>
                  <Text style={styles.drawerUserName}>
                    {currentUser?.username || 'amiri2sb'}
                  </Text>
                  <Text style={styles.drawerUserSub}>
                    {currentUser?.group_name ? `${currentUser.group_name} • subly` : 'subly.app'}
                  </Text>
                </View>
              </View>
              <TouchableOpacity
                onPress={() => closeDrawer()}
                hitSlop={{ top: 12, bottom: 12, left: 12, right: 12 }}
              >
                <Ionicons name="close" size={24} color="#7E8696" />
              </TouchableOpacity>
            </View>

            <View style={styles.drawerDivider} />

            {/* Пункты меню */}
            <View style={styles.drawerNav}>
              <TouchableOpacity
                style={styles.drawerItem}
                activeOpacity={0.7}
                onPress={() => {
                  closeDrawer(() => handleContactDev());
                }}
              >
                <View style={[styles.drawerItemIcon, { backgroundColor: 'rgba(0, 229, 255, 0.12)' }]}>
                  <Ionicons name="paper-plane" size={18} color="#00E5FF" />
                </View>
                <View style={styles.drawerItemTextBox}>
                  <Text style={styles.drawerItemTitle}>Связаться с разработчиком</Text>
                  <Text style={styles.drawerItemSub}>Telegram: @amiri2sb</Text>
                </View>
                <Ionicons name="chevron-forward" size={16} color="#4A5568" />
              </TouchableOpacity>

              <View style={styles.drawerItem}>
                <View style={[styles.drawerItemIcon, { backgroundColor: 'rgba(16, 185, 129, 0.12)' }]}>
                  <Ionicons name="shield-checkmark" size={18} color="#10B981" />
                </View>
                <View style={styles.drawerItemTextBox}>
                  <Text style={styles.drawerItemTitle}>Статус системы</Text>
                  <Text style={styles.drawerItemSub}>MAX WSS: Активен • v2.0.0</Text>
                </View>
              </View>

              <View style={styles.drawerItem}>
                <View style={[styles.drawerItemIcon, { backgroundColor: 'rgba(139, 92, 246, 0.12)' }]}>
                  <Ionicons name="eye-off" size={18} color="#A78BFA" />
                </View>
                <View style={styles.drawerItemTextBox}>
                  <Text style={styles.drawerItemTitle}>Стелс-режим пушей</Text>
                  <Text style={styles.drawerItemSub}>Локскрин: «Есть изменения»</Text>
                </View>
              </View>
            </View>

            {/* Футер с кнопкой выхода в самом низу */}
            <View style={styles.drawerFooter}>
              <TouchableOpacity
                style={styles.logoutButton}
                activeOpacity={0.8}
                onPress={handleLogoutPress}
              >
                <Ionicons name="log-out-outline" size={18} color="#EF4444" />
                <Text style={styles.logoutText}>Выйти из аккаунта</Text>
              </TouchableOpacity>
            </View>
          </Animated.View>
        </View>
      </Modal>

      {/* =================== МОДАЛКА СОЗДАНИЯ ЗАДАЧИ =================== */}
      <Modal
        visible={isModalOpen}
        transparent
        animationType="slide"
        onRequestClose={() => setIsModalOpen(false)}
      >
        <KeyboardAvoidingView
          behavior={Platform.OS === 'ios' ? 'padding' : undefined}
          style={styles.modalOverlay}
        >
          <View style={styles.modalContent}>
            <View style={styles.modalHeader}>
              <Text style={styles.modalTitle}>Новая задача</Text>
              <TouchableOpacity
                onPress={() => setIsModalOpen(false)}
                hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
              >
                <Ionicons name="close" size={22} color="#7E8696" />
              </TouchableOpacity>
            </View>

            {/* Название задачи */}
            <TextInput
              style={styles.modalInput}
              placeholder="Что нужно сделать?"
              placeholderTextColor="#5E6778"
              value={newTitle}
              onChangeText={setNewTitle}
              autoFocus
            />

            {/* Быстрый выбор времени */}
            <Text style={styles.modalSubLabel}>Время / Дедлайн</Text>
            <View style={styles.chipsRow}>
              {QUICK_TIMES.map((time) => {
                const isSelected = selectedTime === time && !customTime;
                return (
                  <TouchableOpacity
                    key={time}
                    style={[styles.chip, isSelected && styles.chipActive]}
                    onPress={() => {
                      setSelectedTime(time);
                      setCustomTime('');
                    }}
                  >
                    <Text style={[styles.chipText, isSelected && styles.chipTextActive]}>
                      {time}
                    </Text>
                  </TouchableOpacity>
                );
              })}
            </View>

            {/* Кастомное время */}
            <TextInput
              style={styles.modalTimeInput}
              placeholder="Или напишите своё время (например 16:30)"
              placeholderTextColor="#5E6778"
              value={customTime}
              onChangeText={setCustomTime}
            />

            {/* Категории */}
            <Text style={styles.modalSubLabel}>Категория</Text>
            <View style={styles.chipsRow}>
              {CATEGORIES.map((cat) => {
                const isSelected = selectedCategory === cat;
                return (
                  <TouchableOpacity
                    key={cat}
                    style={[styles.chip, isSelected && styles.chipActive]}
                    onPress={() => setSelectedCategory(cat)}
                  >
                    <Text style={[styles.chipText, isSelected && styles.chipTextActive]}>
                      {cat}
                    </Text>
                  </TouchableOpacity>
                );
              })}
            </View>

            {/* Тумблер закрепа */}
            <TouchableOpacity
              style={[styles.pinToggleRow, isNewPinned && styles.pinToggleRowActive]}
              activeOpacity={0.8}
              onPress={() => setIsNewPinned(!isNewPinned)}
            >
              <View style={styles.pinToggleLeft}>
                <Ionicons
                  name={isNewPinned ? 'pin' : 'pin-outline'}
                  size={18}
                  color={isNewPinned ? '#00E5FF' : '#7E8696'}
                />
                <Text style={[styles.pinToggleText, isNewPinned && styles.pinToggleTextActive]}>
                  Закрепить вверху (неоновый контур)
                </Text>
              </View>
              <Ionicons
                name={isNewPinned ? 'checkbox' : 'square-outline'}
                size={22}
                color={isNewPinned ? '#00E5FF' : '#5E6778'}
              />
            </TouchableOpacity>

            {/* Кнопки */}
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

      {/* =================== ШТОРКА РАСПИСАНИЯ =================== */}
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
  headerLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
  },
  burgerBtn: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: '#121622',
    borderWidth: 1,
    borderColor: '#1E2536',
    justifyContent: 'center',
    alignItems: 'center',
  },
  headerTitleBox: {
    justifyContent: 'center',
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
    fontSize: 34,
    fontWeight: '800',
    marginTop: 2,
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
  tasksListContainer: {
    gap: 10,
  },
  emptyCard: {
    backgroundColor: '#121622',
    borderRadius: 18,
    paddingVertical: 18,
    borderWidth: 1,
    borderColor: '#1E2536',
    alignItems: 'center',
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
    paddingRight: 8,
  },
  titleWrapper: {
    flex: 1,
  },
  taskTitle: {
    color: '#FFFFFF',
    fontSize: 15,
    fontWeight: '500',
    marginBottom: 4,
  },
  collegeTaskTitle: {
    fontWeight: '700',
    color: '#E0F7FA',
  },
  completedTitle: {
    color: '#636C7E',
    textDecorationLine: 'line-through',
  },
  badgesRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    flexWrap: 'wrap',
  },
  categoryBadge: {
    backgroundColor: '#182133',
    paddingHorizontal: 7,
    paddingVertical: 2,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#22304A',
  },
  categoryBadgeText: {
    color: '#A0AEC0',
    fontSize: 11,
    fontWeight: '600',
  },
  badge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(0, 229, 255, 0.1)',
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: 'rgba(0, 229, 255, 0.25)',
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
  burgerOverlay: {
    flex: 1,
    flexDirection: 'row',
  },
  burgerBackdrop: {
    ...StyleSheet.absoluteFill,
    backgroundColor: 'rgba(0, 0, 0, 0.72)',
  },
  drawerContainer: {
    height: '100%',
    backgroundColor: '#0E131E',
    paddingHorizontal: 18,
    paddingTop: Platform.OS === 'ios' ? 56 : 32,
    paddingBottom: Platform.OS === 'ios' ? 36 : 24,
    borderRightWidth: 1,
    borderRightColor: '#1E2536',
    shadowColor: '#000000',
    shadowOffset: { width: 8, height: 0 },
    shadowOpacity: 0.85,
    shadowRadius: 18,
    elevation: 16,
    justifyContent: 'space-between',
  },
  drawerHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  drawerProfileRow: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  drawerAvatar: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: 'rgba(0, 229, 255, 0.15)',
    borderWidth: 1.5,
    borderColor: '#00E5FF',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  drawerAvatarText: {
    color: '#00E5FF',
    fontSize: 16,
    fontWeight: '800',
  },
  drawerUserInfo: {
    flex: 1,
  },
  drawerUserName: {
    color: '#FFFFFF',
    fontSize: 18,
    fontWeight: '800',
  },
  drawerUserSub: {
    color: '#636C7E',
    fontSize: 12,
    marginTop: 2,
  },
  drawerDivider: {
    height: 1,
    backgroundColor: '#1A2130',
    marginVertical: 18,
  },
  drawerNav: {
    flex: 1,
    gap: 12,
  },
  drawerItem: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#171D2B',
    paddingVertical: 12,
    paddingHorizontal: 14,
    borderRadius: 16,
    borderWidth: 1,
    borderColor: '#222B3D',
  },
  drawerItemIcon: {
    width: 36,
    height: 36,
    borderRadius: 10,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
  },
  drawerItemTextBox: {
    flex: 1,
  },
  drawerItemTitle: {
    color: '#FFFFFF',
    fontSize: 14,
    fontWeight: '600',
  },
  drawerItemSub: {
    color: '#7E8696',
    fontSize: 11,
    marginTop: 2,
  },
  drawerFooter: {
    paddingTop: 16,
    borderTopWidth: 1,
    borderTopColor: '#1A2130',
  },
  logoutButton: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: 'rgba(239, 68, 68, 0.1)',
    borderRadius: 14,
    paddingVertical: 12,
    gap: 8,
    borderWidth: 1,
    borderColor: 'rgba(239, 68, 68, 0.25)',
  },
  logoutText: {
    color: '#EF4444',
    fontSize: 14,
    fontWeight: '700',
  },
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.82)',
    justifyContent: 'flex-end',
  },
  modalContent: {
    backgroundColor: '#121622',
    borderTopLeftRadius: 28,
    borderTopRightRadius: 28,
    paddingHorizontal: 22,
    paddingTop: 20,
    paddingBottom: Platform.OS === 'ios' ? 40 : 24,
    borderWidth: 1,
    borderColor: '#222B3D',
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  modalTitle: {
    color: '#FFFFFF',
    fontSize: 20,
    fontWeight: '800',
  },
  modalSubLabel: {
    color: '#8A93A6',
    fontSize: 12,
    fontWeight: '700',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    marginBottom: 8,
  },
  modalInput: {
    backgroundColor: '#171D2B',
    borderRadius: 14,
    paddingHorizontal: 16,
    height: 50,
    color: '#FFFFFF',
    fontSize: 16,
    borderWidth: 1,
    borderColor: '#263145',
    marginBottom: 16,
  },
  modalTimeInput: {
    backgroundColor: '#171D2B',
    borderRadius: 12,
    paddingHorizontal: 14,
    height: 42,
    color: '#FFFFFF',
    fontSize: 13,
    borderWidth: 1,
    borderColor: '#263145',
    marginBottom: 16,
  },
  chipsRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
    marginBottom: 10,
  },
  chip: {
    backgroundColor: '#182133',
    paddingHorizontal: 12,
    paddingVertical: 7,
    borderRadius: 20,
    borderWidth: 1,
    borderColor: '#22304A',
  },
  chipActive: {
    backgroundColor: 'rgba(0, 229, 255, 0.15)',
    borderColor: '#00E5FF',
  },
  chipText: {
    color: '#8A93A6',
    fontSize: 12,
    fontWeight: '600',
  },
  chipTextActive: {
    color: '#00E5FF',
    fontWeight: '700',
  },
  pinToggleRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    backgroundColor: '#171D2B',
    borderRadius: 14,
    paddingHorizontal: 14,
    paddingVertical: 12,
    borderWidth: 1,
    borderColor: '#263145',
    marginBottom: 20,
  },
  pinToggleRowActive: {
    borderColor: '#00E5FF',
    backgroundColor: 'rgba(0, 229, 255, 0.08)',
  },
  pinToggleLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 10,
  },
  pinToggleText: {
    color: '#8A93A6',
    fontSize: 14,
    fontWeight: '600',
  },
  pinToggleTextActive: {
    color: '#00E5FF',
    fontWeight: '700',
  },
  modalButtons: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    alignItems: 'center',
    gap: 12,
  },
  modalCancel: {
    paddingVertical: 12,
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
    paddingVertical: 12,
    paddingHorizontal: 26,
    borderRadius: 14,
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.5,
    shadowRadius: 8,
  },
  modalSubmitText: {
    color: '#000000',
    fontSize: 15,
    fontWeight: '800',
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
    borderColor: '#263145',
  },
  sheetSlotRep: {
    borderColor: '#00E5FF',
    backgroundColor: 'rgba(0, 229, 255, 0.08)',
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
    backgroundColor: '#00E5FF',
    borderRadius: 6,
    paddingHorizontal: 6,
    paddingVertical: 2,
  },
  repPillText: {
    color: '#000000',
    fontSize: 10,
    fontWeight: '800',
  },
  sheetSubject: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '700',
    marginBottom: 6,
  },
  sheetSlotBottom: {
    flexDirection: 'row',
    gap: 16,
  },
  sheetMeta: {
    color: '#8A93A6',
    fontSize: 13,
  },
  sheetFullBtn: {
    backgroundColor: '#00E5FF',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 14,
    borderRadius: 16,
    gap: 8,
  },
  sheetFullBtnText: {
    color: '#000000',
    fontSize: 15,
    fontWeight: '700',
  },
});
