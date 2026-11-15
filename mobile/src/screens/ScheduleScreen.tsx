import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  View,
  Text,
  ScrollView,
  TouchableOpacity,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { ScheduleSlot, Task } from '../types';
import { fetchSchedule, fetchTasks, updateTaskStatus } from '../api/client';

export const ScheduleScreen: React.FC = () => {
  const todayDayNum = 15; // 15 Nov 2026
  const [selectedDayNum, setSelectedDayNum] = useState<number>(todayDayNum);
  const [activeSegment, setActiveSegment] = useState<'tasks' | 'schedule'>('tasks');
  const [tasks, setTasks] = useState<Task[]>([]);
  const [schedule, setSchedule] = useState<ScheduleSlot[]>([]);

  // Дни ноября 2026: 1 ноября 2026 - Воскресенье (0)
  // Всего 30 дней в ноябре
  const daysInMonth = 30;
  const monthName = 'Ноябрь 2026';
  const weekDays = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'];

  useEffect(() => {
    loadData();
  }, [selectedDayNum]);

  // Вычисление дня недели для выбранного числа ноября 2026 (1 нояб = Вс = 0)
  const getDayOfWeek = (dayNum: number): number => {
    return (dayNum - 1) % 7; // 0 - Вс, 1 - Пн ... 6 - Сб
  };

  const getFormattedDate = (dayNum: number): string => {
    const dayStr = dayNum < 10 ? `0${dayNum}` : `${dayNum}`;
    return `2026-11-${dayStr}`;
  };

  const loadData = async () => {
    const dayOfWeek = getDayOfWeek(selectedDayNum);
    const [taskList, scheduleList] = await Promise.all([
      fetchTasks(),
      fetchSchedule(dayOfWeek),
    ]);
    setTasks(taskList);
    setSchedule(scheduleList);
  };

  const toggleTask = (task: Task) => {
    const newStatus = task.status === 'completed' ? 'pending' : 'completed';
    setTasks((prev) =>
      prev.map((t) => (t.id === task.id ? { ...t, status: newStatus } : t))
    );
    updateTaskStatus(task.id, newStatus);
  };

  const selectedDateStr = getFormattedDate(selectedDayNum);

  // Фильтрация задач для выбранного дня
  const dayTasks = tasks.filter((t) => {
    if (t.status === 'completed') {
      return t.completed_at === selectedDateStr || t.date === selectedDateStr;
    }
    return t.date === selectedDateStr;
  });

  const pendingDayTasks = dayTasks.filter((t) => t.status !== 'completed');
  const completedDayTasks = dayTasks.filter((t) => t.status === 'completed');

  const selectedDayOfWeek = getDayOfWeek(selectedDayNum);
  const isSunday = selectedDayOfWeek === 0;
  const replacementSlot = schedule.find((s) => s.is_replacement);

  // Проверка, есть ли задачи или история для конкретного дня месяца
  const hasTasksOnDay = (dayNum: number): boolean => {
    const dateStr = getFormattedDate(dayNum);
    return tasks.some(
      (t) => t.date === dateStr || t.completed_at === dateStr
    );
  };

  return (
    <View style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        {/* Заголовок месяца */}
        <View style={styles.header}>
          <View>
            <Text style={styles.headerCategory}>Subly / Календарь</Text>
            <Text style={styles.monthTitle}>{monthName}</Text>
          </View>
          <View style={styles.headerRight}>
            <TouchableOpacity
              style={[
                styles.todayButton,
                selectedDayNum === todayDayNum && styles.todayButtonActive,
              ]}
              onPress={() => setSelectedDayNum(todayDayNum)}
            >
              <Text
                style={[
                  styles.todayButtonText,
                  selectedDayNum === todayDayNum && styles.todayButtonTextActive,
                ]}
              >
                Сегодня
              </Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Календарная сетка месяца */}
        <View style={styles.calendarCard}>
          {/* Дни недели */}
          <View style={styles.weekDaysRow}>
            {weekDays.map((wd, idx) => (
              <Text
                key={wd}
                style={[styles.weekDayText, (idx === 0 || idx === 6) && styles.weekendText]}
              >
                {wd}
              </Text>
            ))}
          </View>

          {/* Сетка чисел (30 дней ноября) */}
          <View style={styles.daysGrid}>
            {Array.from({ length: daysInMonth }, (_, i) => i + 1).map((day) => {
              const isToday = day === todayDayNum;
              const isSelected = day === selectedDayNum;
              const hasTasks = hasTasksOnDay(day);
              const dayOfWeek = getDayOfWeek(day);
              const isWeekend = dayOfWeek === 0 || dayOfWeek === 6;

              return (
                <TouchableOpacity
                  key={day}
                  style={[
                    styles.dayCell,
                    isSelected && styles.dayCellSelected,
                    isToday && !isSelected && styles.dayCellToday,
                  ]}
                  activeOpacity={0.7}
                  onPress={() => setSelectedDayNum(day)}
                >
                  <Text
                    style={[
                      styles.dayCellText,
                      isWeekend && styles.dayWeekendText,
                      isToday && styles.dayTodayText,
                      isSelected && styles.daySelectedText,
                    ]}
                  >
                    {day}
                  </Text>

                  {/* Индикатор задач на день */}
                  <View style={styles.dotsRow}>
                    {hasTasks && (
                      <View
                        style={[
                          styles.taskDot,
                          isSelected && styles.taskDotSelected,
                          isToday && !isSelected && styles.taskDotToday,
                        ]}
                      />
                    )}
                  </View>
                </TouchableOpacity>
              );
            })}
          </View>
        </View>

        {/* Переключатель: История задач / Расписание пар */}
        <View style={styles.segmentContainer}>
          <TouchableOpacity
            style={[
              styles.segmentButton,
              activeSegment === 'tasks' && styles.segmentButtonActive,
            ]}
            onPress={() => setActiveSegment('tasks')}
            activeOpacity={0.8}
          >
            <Ionicons
              name="checkbox"
              size={16}
              color={activeSegment === 'tasks' ? '#00E5FF' : '#7E8696'}
            />
            <Text
              style={[
                styles.segmentText,
                activeSegment === 'tasks' && styles.segmentTextActive,
              ]}
            >
              История задач ({dayTasks.length})
            </Text>
          </TouchableOpacity>

          <TouchableOpacity
            style={[
              styles.segmentButton,
              activeSegment === 'schedule' && styles.segmentButtonActive,
            ]}
            onPress={() => setActiveSegment('schedule')}
            activeOpacity={0.8}
          >
            <Ionicons
              name="school"
              size={16}
              color={activeSegment === 'schedule' ? '#00E5FF' : '#7E8696'}
            />
            <Text
              style={[
                styles.segmentText,
                activeSegment === 'schedule' && styles.segmentTextActive,
              ]}
            >
              Пары в колледже {isSunday ? '(Вых)' : `(${schedule.length})`}
            </Text>
          </TouchableOpacity>
        </View>

        {/* Секция: История и задачи дня */}
        {activeSegment === 'tasks' && (
          <View style={styles.contentSection}>
            {dayTasks.length === 0 ? (
              <View style={styles.emptyCard}>
                <Ionicons name="calendar-outline" size={38} color="#353D50" />
                <Text style={styles.emptyTitle}>Нет задач на этот день</Text>
                <Text style={styles.emptySubtitle}>
                  {selectedDayNum} ноября — задач не запланировано и не выполнялось
                </Text>
              </View>
            ) : (
              <>
                {pendingDayTasks.length > 0 && (
                  <View style={styles.taskGroup}>
                    <Text style={styles.groupHeaderTitle}>
                      В процессе ({pendingDayTasks.length})
                    </Text>
                    {pendingDayTasks.map((t) => (
                      <View key={t.id} style={styles.taskItemCard}>
                        <TouchableOpacity
                          style={styles.taskCheck}
                          onPress={() => toggleTask(t)}
                        >
                          <Ionicons name="square-outline" size={20} color="#7E8696" />
                        </TouchableOpacity>
                        <View style={styles.taskDetails}>
                          <Text style={styles.taskItemTitle}>{t.title}</Text>
                          <View style={styles.taskMetaRow}>
                            {t.category && (
                              <View style={styles.metaBadge}>
                                <Text style={styles.metaBadgeText}>{t.category}</Text>
                              </View>
                            )}
                            {t.deadline && (
                              <View style={styles.timeBadge}>
                                <Ionicons name="time-outline" size={11} color="#00E5FF" style={{ marginRight: 3 }} />
                                <Text style={styles.timeBadgeText}>{t.deadline}</Text>
                              </View>
                            )}
                          </View>
                        </View>
                      </View>
                    ))}
                  </View>
                )}

                {completedDayTasks.length > 0 && (
                  <View style={styles.taskGroup}>
                    <Text style={styles.groupHeaderTitle}>
                      Выполнено ({completedDayTasks.length})
                    </Text>
                    {completedDayTasks.map((t) => (
                      <View key={t.id} style={styles.taskItemCard}>
                        <TouchableOpacity
                          style={styles.taskCheck}
                          onPress={() => toggleTask(t)}
                        >
                          <Ionicons name="checkbox" size={20} color="#00E5FF" />
                        </TouchableOpacity>
                        <View style={styles.taskDetails}>
                          <Text style={[styles.taskItemTitle, styles.taskCompletedTitle]}>
                            {t.title}
                          </Text>
                          <View style={styles.taskMetaRow}>
                            {t.category && (
                              <View style={styles.metaBadge}>
                                <Text style={styles.metaBadgeText}>{t.category}</Text>
                              </View>
                            )}
                            <View style={styles.doneBadge}>
                              <Ionicons name="checkmark" size={11} color="#10B981" style={{ marginRight: 3 }} />
                              <Text style={styles.doneBadgeText}>Завершено</Text>
                            </View>
                          </View>
                        </View>
                      </View>
                    ))}
                  </View>
                )}
              </>
            )}
          </View>
        )}

        {/* Секция: Пары в колледже */}
        {activeSegment === 'schedule' && (
          <View style={styles.contentSection}>
            {isSunday ? (
              <View style={styles.emptyCard}>
                <Ionicons name="sunny-outline" size={42} color="#00E5FF" />
                <Text style={styles.emptyTitle}>Воскресенье — выходной!</Text>
                <Text style={styles.emptySubtitle}>
                  Пар нет. Отдыхаем, кодим свои пет-проекты или готовимся к новой неделе.
                </Text>
              </View>
            ) : (
              <>
                {replacementSlot && (
                  <View style={styles.alertBanner}>
                    <View style={styles.alertHeader}>
                      <Ionicons name="swap-horizontal" size={20} color="#00E5FF" />
                      <Text style={styles.alertTitle}>Замена в расписании (из MAX)</Text>
                    </View>
                    <Text style={styles.alertText}>
                      {replacementSlot.order_num} пара — {replacementSlot.subject} ({replacementSlot.room})
                    </Text>
                  </View>
                )}

                <View style={styles.timeline}>
                  {schedule.map((slot) => {
                    const isReplacement = slot.is_replacement;
                    return (
                      <View
                        key={slot.id}
                        style={[
                          styles.slotCard,
                          isReplacement && styles.replacementSlotCard,
                        ]}
                      >
                        <View style={styles.slotHeader}>
                          <Text style={[styles.orderText, isReplacement && styles.repOrderText]}>
                            {slot.order_num} пара ({slot.start_time} - {slot.end_time})
                          </Text>
                          {isReplacement && (
                            <View style={styles.replacementBadge}>
                              <Text style={styles.replacementBadgeText}>Замена</Text>
                            </View>
                          )}
                        </View>

                        <Text style={styles.subjectText}>{slot.subject}</Text>

                        <View style={styles.slotFooter}>
                          <View style={styles.slotMeta}>
                            <Ionicons name="location-outline" size={14} color="#8A93A6" />
                            <Text style={styles.metaText}>{slot.room}</Text>
                          </View>
                          {slot.teacher && (
                            <View style={styles.slotMeta}>
                              <Ionicons name="person-outline" size={14} color="#8A93A6" />
                              <Text style={styles.metaText}>{slot.teacher}</Text>
                            </View>
                          )}
                        </View>
                      </View>
                    );
                  })}
                </View>
              </>
            )}
          </View>
        )}
      </ScrollView>
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
    marginBottom: 16,
  },
  headerCategory: {
    color: '#7E8696',
    fontSize: 13,
    fontWeight: '600',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  monthTitle: {
    color: '#FFFFFF',
    fontSize: 32,
    fontWeight: '800',
    marginTop: 4,
    letterSpacing: -0.5,
  },
  headerRight: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  todayButton: {
    backgroundColor: '#121622',
    paddingHorizontal: 14,
    paddingVertical: 7,
    borderRadius: 20,
    borderWidth: 1,
    borderColor: '#202738',
  },
  todayButtonActive: {
    borderColor: '#00E5FF',
    backgroundColor: 'rgba(0, 229, 255, 0.1)',
  },
  todayButtonText: {
    color: '#8A93A6',
    fontSize: 13,
    fontWeight: '600',
  },
  todayButtonTextActive: {
    color: '#00E5FF',
    fontWeight: '700',
  },
  calendarCard: {
    backgroundColor: '#121622',
    borderRadius: 24,
    padding: 16,
    borderWidth: 1,
    borderColor: '#1E2536',
    marginBottom: 16,
  },
  weekDaysRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 12,
    paddingBottom: 8,
    borderBottomWidth: 1,
    borderBottomColor: '#1A2130',
  },
  weekDayText: {
    color: '#7E8696',
    fontSize: 12,
    fontWeight: '700',
    width: 36,
    textAlign: 'center',
  },
  weekendText: {
    color: '#5E6778',
  },
  daysGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-around',
    gap: 4,
  },
  dayCell: {
    width: 40,
    height: 44,
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 14,
  },
  dayCellToday: {
    borderWidth: 1.5,
    borderColor: '#FF7A00',
    backgroundColor: 'rgba(255, 122, 0, 0.12)',
  },
  dayCellSelected: {
    backgroundColor: '#00E5FF',
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.6,
    shadowRadius: 8,
    elevation: 4,
  },
  dayCellText: {
    color: '#FFFFFF',
    fontSize: 14,
    fontWeight: '600',
  },
  dayWeekendText: {
    color: '#8A93A6',
  },
  dayTodayText: {
    color: '#FF9100',
    fontWeight: '800',
  },
  daySelectedText: {
    color: '#000000',
    fontWeight: '800',
  },
  dotsRow: {
    height: 6,
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 2,
  },
  taskDot: {
    width: 4,
    height: 4,
    borderRadius: 2,
    backgroundColor: '#00E5FF',
  },
  taskDotSelected: {
    backgroundColor: '#000000',
  },
  taskDotToday: {
    backgroundColor: '#FF7A00',
  },
  segmentContainer: {
    flexDirection: 'row',
    backgroundColor: '#121622',
    borderRadius: 20,
    padding: 4,
    borderWidth: 1,
    borderColor: '#1E2536',
    marginBottom: 16,
  },
  segmentButton: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 10,
    borderRadius: 16,
    gap: 6,
  },
  segmentButtonActive: {
    backgroundColor: '#1A2333',
    borderWidth: 1,
    borderColor: 'rgba(0, 229, 255, 0.3)',
  },
  segmentText: {
    color: '#7E8696',
    fontSize: 12,
    fontWeight: '600',
  },
  segmentTextActive: {
    color: '#00E5FF',
    fontWeight: '700',
  },
  contentSection: {
    gap: 12,
  },
  emptyCard: {
    backgroundColor: '#121622',
    borderRadius: 20,
    padding: 24,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: '#1E2536',
  },
  emptyTitle: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '700',
    marginTop: 12,
  },
  emptySubtitle: {
    color: '#7E8696',
    fontSize: 13,
    textAlign: 'center',
    marginTop: 6,
    lineHeight: 18,
  },
  taskGroup: {
    marginBottom: 8,
  },
  groupHeaderTitle: {
    color: '#CAD1DE',
    fontSize: 14,
    fontWeight: '700',
    marginBottom: 8,
    paddingHorizontal: 4,
  },
  taskItemCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#121622',
    borderRadius: 16,
    paddingVertical: 12,
    paddingHorizontal: 14,
    borderWidth: 1,
    borderColor: '#1E2536',
    marginBottom: 8,
  },
  taskCheck: {
    marginRight: 12,
  },
  taskDetails: {
    flex: 1,
  },
  taskItemTitle: {
    color: '#FFFFFF',
    fontSize: 14,
    fontWeight: '500',
    marginBottom: 4,
  },
  taskCompletedTitle: {
    color: '#636C7E',
    textDecorationLine: 'line-through',
  },
  taskMetaRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
  },
  metaBadge: {
    backgroundColor: '#182133',
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 6,
    borderWidth: 1,
    borderColor: '#22304A',
  },
  metaBadgeText: {
    color: '#A0AEC0',
    fontSize: 10,
    fontWeight: '600',
  },
  timeBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(0, 229, 255, 0.1)',
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 6,
    borderWidth: 1,
    borderColor: 'rgba(0, 229, 255, 0.25)',
  },
  timeBadgeText: {
    color: '#00E5FF',
    fontSize: 10,
    fontWeight: '600',
  },
  doneBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: 'rgba(16, 185, 129, 0.12)',
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 6,
  },
  doneBadgeText: {
    color: '#10B981',
    fontSize: 10,
    fontWeight: '600',
  },
  alertBanner: {
    backgroundColor: '#0F1A24',
    borderWidth: 1,
    borderColor: '#00E5FF',
    borderRadius: 16,
    padding: 14,
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.35,
    shadowRadius: 8,
    elevation: 4,
    marginBottom: 12,
  },
  alertHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 4,
  },
  alertTitle: {
    color: '#00E5FF',
    fontSize: 14,
    fontWeight: '700',
    marginLeft: 8,
  },
  alertText: {
    color: '#D1E8FF',
    fontSize: 13,
  },
  timeline: {
    gap: 10,
  },
  slotCard: {
    backgroundColor: '#121622',
    borderRadius: 18,
    padding: 16,
    borderWidth: 1,
    borderColor: '#1E2536',
  },
  replacementSlotCard: {
    borderColor: '#00E5FF',
    borderWidth: 1.5,
    backgroundColor: '#0F1826',
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.3,
    shadowRadius: 8,
  },
  slotHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  orderText: {
    color: '#7E8696',
    fontSize: 12,
    fontWeight: '600',
  },
  repOrderText: {
    color: '#00E5FF',
    fontWeight: '700',
  },
  replacementBadge: {
    backgroundColor: '#00E5FF',
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 8,
  },
  replacementBadgeText: {
    color: '#000000',
    fontSize: 10,
    fontWeight: '800',
  },
  subjectText: {
    color: '#FFFFFF',
    fontSize: 17,
    fontWeight: '700',
    marginBottom: 10,
  },
  slotFooter: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 16,
  },
  slotMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  metaText: {
    color: '#8A93A6',
    fontSize: 13,
  },
});
