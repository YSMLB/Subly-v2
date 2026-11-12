import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  Alert,
  Platform,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { ProposedChange } from '../types';
import { fetchPendingChanges, reviewChange } from '../api/client';

export const ModerationScreen: React.FC = () => {
  const [changes, setChanges] = useState<ProposedChange[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadChanges();
  }, []);

  const loadChanges = async () => {
    setLoading(true);
    const data = await fetchPendingChanges();
    setChanges(data);
    setLoading(false);
  };

  const handleApprove = async (change: ProposedChange) => {
    await reviewChange(change.id, 'approved');
    setChanges((prev) => prev.filter((c) => c.id !== change.id));
    Alert.alert('Замена одобрена', 'Расписание группы успешно обновлено в базе данных.');
  };

  const handleReject = async (change: ProposedChange) => {
    await reviewChange(change.id, 'rejected');
    setChanges((prev) => prev.filter((c) => c.id !== change.id));
    Alert.alert('Замена отклонена', 'Заявка помечена как отклоненная.');
  };

  return (
    <View style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        <View style={styles.header}>
          <Text style={styles.headerCategory}>Subly / Панель старосты</Text>
          <Text style={styles.headerTitle}>Модерация</Text>
          <Text style={styles.headerSubtitle}>
            Входящие замены из MAX Messenger, ожидающие подтверждения
          </Text>
        </View>

        {changes.length === 0 ? (
          <View style={styles.emptyCard}>
            <Ionicons name="shield-checkmark" size={48} color="#00E5FF" />
            <Text style={styles.emptyTitle}>Всё чисто</Text>
            <Text style={styles.emptyText}>
              Новых замен от старосты пока нет. Когда в чате MAX появится сообщение, оно отобразится здесь автоматически.
            </Text>
            <TouchableOpacity style={styles.refreshButton} onPress={loadChanges}>
              <Ionicons name="refresh" size={16} color="#000000" />
              <Text style={styles.refreshText}>Обновить</Text>
            </TouchableOpacity>
          </View>
        ) : (
          changes.map((item) => (
            <View key={item.id} style={styles.changeCard}>
              <View style={styles.cardHeader}>
                <View style={styles.groupBadge}>
                  <Text style={styles.groupText}>{item.group_name || 'Группа'}</Text>
                </View>
                <Text style={styles.dateText}>{item.created_at || 'Сегодня'}</Text>
              </View>

              <Text style={styles.rawTitle}>Текст из чата:</Text>
              <View style={styles.rawBox}>
                <Text style={styles.rawText}>{item.raw_text}</Text>
              </View>

              <View style={styles.actions}>
                <TouchableOpacity
                  style={[styles.actionBtn, styles.rejectBtn]}
                  onPress={() => handleReject(item)}
                  activeOpacity={0.8}
                >
                  <Ionicons name="close" size={18} color="#EF4444" />
                  <Text style={styles.rejectText}>Отклонить</Text>
                </TouchableOpacity>

                <TouchableOpacity
                  style={[styles.actionBtn, styles.approveBtn]}
                  onPress={() => handleApprove(item)}
                  activeOpacity={0.8}
                >
                  <Ionicons name="checkmark" size={18} color="#000000" />
                  <Text style={styles.approveText}>Одобрить</Text>
                </TouchableOpacity>
              </View>
            </View>
          ))
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
    marginBottom: 20,
  },
  headerCategory: {
    color: '#00E5FF',
    fontSize: 13,
    fontWeight: '700',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  headerTitle: {
    color: '#FFFFFF',
    fontSize: 34,
    fontWeight: '800',
    marginTop: 4,
    letterSpacing: -0.5,
  },
  headerSubtitle: {
    color: '#7E8696',
    fontSize: 14,
    marginTop: 6,
    lineHeight: 20,
  },
  emptyCard: {
    backgroundColor: '#121622',
    borderRadius: 24,
    padding: 32,
    alignItems: 'center',
    borderWidth: 1,
    borderColor: '#1E2536',
    marginTop: 40,
  },
  emptyTitle: {
    color: '#FFFFFF',
    fontSize: 20,
    fontWeight: '700',
    marginTop: 16,
  },
  emptyText: {
    color: '#7E8696',
    fontSize: 14,
    textAlign: 'center',
    marginTop: 8,
    lineHeight: 20,
  },
  refreshButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#00E5FF',
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 14,
    marginTop: 20,
    gap: 8,
  },
  refreshText: {
    color: '#000000',
    fontSize: 14,
    fontWeight: '700',
  },
  changeCard: {
    backgroundColor: '#121622',
    borderRadius: 22,
    padding: 18,
    borderWidth: 1,
    borderColor: '#1E2536',
    marginBottom: 16,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  groupBadge: {
    backgroundColor: 'rgba(0, 229, 255, 0.15)',
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 10,
    borderWidth: 1,
    borderColor: 'rgba(0, 229, 255, 0.3)',
  },
  groupText: {
    color: '#00E5FF',
    fontSize: 12,
    fontWeight: '700',
  },
  dateText: {
    color: '#7E8696',
    fontSize: 13,
  },
  rawTitle: {
    color: '#8A93A6',
    fontSize: 13,
    fontWeight: '600',
    marginBottom: 6,
  },
  rawBox: {
    backgroundColor: '#171D2B',
    borderRadius: 14,
    padding: 14,
    borderWidth: 1,
    borderColor: '#242D3F',
    marginBottom: 16,
  },
  rawText: {
    color: '#CAD1DE',
    fontSize: 14,
    lineHeight: 20,
    fontFamily: Platform.OS === 'ios' ? 'Menlo' : 'monospace',
  },
  actions: {
    flexDirection: 'row',
    gap: 12,
  },
  actionBtn: {
    flex: 1,
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    height: 48,
    borderRadius: 14,
    gap: 6,
  },
  rejectBtn: {
    backgroundColor: 'rgba(239, 68, 68, 0.1)',
    borderWidth: 1,
    borderColor: 'rgba(239, 68, 68, 0.3)',
  },
  rejectText: {
    color: '#EF4444',
    fontSize: 14,
    fontWeight: '700',
  },
  approveBtn: {
    backgroundColor: '#00E5FF',
  },
  approveText: {
    color: '#000000',
    fontSize: 14,
    fontWeight: '700',
  },
});
