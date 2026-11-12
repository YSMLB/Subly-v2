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
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { Task } from '../types';
import { fetchTasks, updateTaskStatus, createTask } from '../api/client';

export const TasksScreen: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newTitle, setNewTitle] = useState('');

  useEffect(() => {
    loadTasks();
  }, []);

  const loadTasks = async () => {
    const list = await fetchTasks();
    setTasks(list);
  };

  const toggleTask = (task: Task) => {
    const newStatus = task.status === 'completed' ? 'pending' : 'completed';
    setTasks((prev) =>
      prev.map((t) => (t.id === task.id ? { ...t, status: newStatus } : t))
    );
    updateTaskStatus(task.id, newStatus);
  };

  const handleCreateTask = async () => {
    if (!newTitle.trim()) return;
    const created = await createTask(1, newTitle.trim());
    setTasks((prev) => [created, ...prev]);
    setNewTitle('');
    setIsModalOpen(false);
  };

  const pinnedTasks = tasks.filter((t) => t.is_pinned && t.status !== 'completed');
  const todayTasks = tasks.filter((t) => !t.is_pinned && t.status !== 'completed');
  const completedTasks = tasks.filter((t) => t.status === 'completed');

  return (
    <View style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        <View style={styles.header}>
          <Text style={styles.headerCategory}>Subly / Задачи</Text>
          <Text style={styles.headerTitle}>Today</Text>
        </View>

        {pinnedTasks.length > 0 && (
          <View style={styles.section}>
            <View style={styles.sectionHeader}>
              <Ionicons name="pin" size={16} color="#00E5FF" />
              <Text style={styles.sectionTitle}>Pinned</Text>
            </View>
            <View style={styles.card}>
              {pinnedTasks.map((task) => (
                <View key={task.id} style={styles.taskRow}>
                  <TouchableOpacity
                    style={styles.checkbox}
                    onPress={() => toggleTask(task)}
                  >
                    <Ionicons name="square-outline" size={22} color="#7E8696" />
                  </TouchableOpacity>
                  <Text style={styles.taskTitle}>{task.title}</Text>
                  <View style={styles.badge}>
                    <Text style={styles.badgeText}>{task.deadline || 'Today'}</Text>
                  </View>
                </View>
              ))}
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
              <Text style={styles.emptyText}>Все задачи на сегодня выполнены 🎉</Text>
            ) : (
              todayTasks.map((task, idx) => (
                <View
                  key={task.id}
                  style={[
                    styles.taskRow,
                    idx !== todayTasks.length - 1 && styles.taskRowBorder,
                  ]}
                >
                  <TouchableOpacity
                    style={styles.checkbox}
                    onPress={() => toggleTask(task)}
                  >
                    <Ionicons name="square-outline" size={22} color="#7E8696" />
                  </TouchableOpacity>
                  <Text style={styles.taskTitle}>{task.title}</Text>
                  {task.deadline && (
                    <View style={styles.badge}>
                      <Text style={styles.badgeText}>{task.deadline}</Text>
                    </View>
                  )}
                </View>
              ))
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
              {completedTasks.map((task, idx) => (
                <View
                  key={task.id}
                  style={[
                    styles.taskRow,
                    idx !== completedTasks.length - 1 && styles.taskRowBorder,
                  ]}
                >
                  <TouchableOpacity
                    style={styles.checkbox}
                    onPress={() => toggleTask(task)}
                  >
                    <Ionicons name="checkbox" size={22} color="#00E5FF" />
                  </TouchableOpacity>
                  <Text style={[styles.taskTitle, styles.completedTitle]}>
                    {task.title}
                  </Text>
                </View>
              ))}
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
              placeholder="Название задачи..."
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
  card: {
    backgroundColor: '#121622',
    borderRadius: 22,
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderWidth: 1,
    borderColor: '#1E2536',
  },
  taskRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 14,
  },
  taskRowBorder: {
    borderBottomWidth: 1,
    borderBottomColor: '#1A2130',
  },
  checkbox: {
    marginRight: 12,
  },
  taskTitle: {
    flex: 1,
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '500',
  },
  completedTitle: {
    color: '#636C7E',
    textDecorationLine: 'line-through',
  },
  badge: {
    backgroundColor: 'rgba(0, 229, 255, 0.12)',
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 12,
    borderWidth: 1,
    borderColor: 'rgba(0, 229, 255, 0.3)',
  },
  badgeText: {
    color: '#00E5FF',
    fontSize: 12,
    fontWeight: '600',
  },
  emptyText: {
    color: '#7E8696',
    fontSize: 14,
    textAlign: 'center',
    paddingVertical: 16,
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
    shadowOpacity: 0.6,
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
});
