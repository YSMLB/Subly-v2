import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  View,
  Text,
  ScrollView,
  TouchableOpacity,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { ScheduleSlot } from '../types';
import { fetchSchedule } from '../api/client';

export const ScheduleScreen: React.FC = () => {
  const [schedule, setSchedule] = useState<ScheduleSlot[]>([]);
  const [selectedDay, setSelectedDay] = useState(3);

  const days = [
    { dayNum: 9, dayName: 'Пн', dayOfWeek: 1 },
    { dayNum: 10, dayName: 'Вт', dayOfWeek: 2 },
    { dayNum: 11, dayName: 'Ср', dayOfWeek: 3 },
    { dayNum: 12, dayName: 'Чт', dayOfWeek: 4 },
    { dayNum: 13, dayName: 'Пт', dayOfWeek: 5 },
    { dayNum: 14, dayName: 'Сб', dayOfWeek: 6 },
  ];

  useEffect(() => {
    loadSchedule();
  }, [selectedDay]);

  const loadSchedule = async () => {
    const data = await fetchSchedule();
    setSchedule(data);
  };

  const replacementSlot = schedule.find((s) => s.is_replacement);

  return (
    <View style={styles.container}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
      >
        <View style={styles.header}>
          <Text style={styles.monthTitle}>Ноябрь</Text>
          <View style={styles.weekStrip}>
            {days.map((item) => {
              const isSelected = item.dayOfWeek === selectedDay;
              return (
                <TouchableOpacity
                  key={item.dayNum}
                  style={[styles.dayItem, isSelected && styles.activeDayItem]}
                  onPress={() => setSelectedDay(item.dayOfWeek)}
                  activeOpacity={0.7}
                >
                  <Text style={[styles.dayName, isSelected && styles.activeDayText]}>
                    {item.dayName}
                  </Text>
                  <Text style={[styles.dayNum, isSelected && styles.activeDayText]}>
                    {item.dayNum}
                  </Text>
                </TouchableOpacity>
              );
            })}
          </View>
        </View>

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
            const isActive = slot.status === 'active';
            const isReplacement = slot.is_replacement;

            return (
              <View
                key={slot.id}
                style={[
                  styles.slotCard,
                  isActive && styles.activeSlotCard,
                  isReplacement && styles.replacementSlotCard,
                ]}
              >
                <View style={styles.slotHeader}>
                  <Text style={[styles.orderText, isActive && styles.activeOrderText]}>
                    {slot.order_num} пара ({slot.start_time} - {slot.end_time})
                  </Text>
                  {isReplacement && (
                    <View style={styles.replacementBadge}>
                      <Text style={styles.replacementBadgeText}>Замена</Text>
                    </View>
                  )}
                  {isActive && !isReplacement && (
                    <View style={styles.activeBadge}>
                      <Text style={styles.activeBadgeText}>Идёт сейчас</Text>
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
  monthTitle: {
    color: '#FFFFFF',
    fontSize: 32,
    fontWeight: '800',
    marginBottom: 16,
    letterSpacing: -0.5,
  },
  weekStrip: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    backgroundColor: '#10131C',
    borderRadius: 20,
    padding: 6,
    borderWidth: 1,
    borderColor: '#1C2333',
  },
  dayItem: {
    alignItems: 'center',
    paddingVertical: 10,
    paddingHorizontal: 12,
    borderRadius: 14,
  },
  activeDayItem: {
    backgroundColor: '#00E5FF',
  },
  dayName: {
    color: '#6B7280',
    fontSize: 12,
    fontWeight: '600',
    marginBottom: 4,
  },
  dayNum: {
    color: '#FFFFFF',
    fontSize: 16,
    fontWeight: '700',
  },
  activeDayText: {
    color: '#000000',
  },
  alertBanner: {
    backgroundColor: '#0F1A26',
    borderRadius: 20,
    padding: 16,
    borderWidth: 1.5,
    borderColor: '#00E5FF',
    marginBottom: 20,
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.5,
    shadowRadius: 10,
    elevation: 6,
  },
  alertHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 6,
  },
  alertTitle: {
    color: '#00E5FF',
    fontSize: 14,
    fontWeight: '700',
    marginLeft: 8,
  },
  alertText: {
    color: '#FFFFFF',
    fontSize: 15,
    fontWeight: '600',
  },
  timeline: {
    gap: 14,
  },
  slotCard: {
    backgroundColor: '#121622',
    borderRadius: 20,
    padding: 18,
    borderWidth: 1,
    borderColor: '#1E2536',
  },
  activeSlotCard: {
    borderColor: '#00E5FF',
    borderWidth: 1.5,
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.4,
    shadowRadius: 8,
    elevation: 5,
  },
  replacementSlotCard: {
    backgroundColor: '#121824',
  },
  slotHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  orderText: {
    color: '#8A93A6',
    fontSize: 13,
    fontWeight: '600',
  },
  activeOrderText: {
    color: '#00E5FF',
  },
  subjectText: {
    color: '#FFFFFF',
    fontSize: 18,
    fontWeight: '700',
    marginBottom: 12,
  },
  slotFooter: {
    flexDirection: 'row',
    gap: 16,
  },
  slotMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
  },
  metaText: {
    color: '#8A93A6',
    fontSize: 13,
  },
  replacementBadge: {
    backgroundColor: 'rgba(0, 229, 255, 0.15)',
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: 'rgba(0, 229, 255, 0.4)',
  },
  replacementBadgeText: {
    color: '#00E5FF',
    fontSize: 11,
    fontWeight: '700',
  },
  activeBadge: {
    backgroundColor: '#1F293D',
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 8,
  },
  activeBadgeText: {
    color: '#CAD1DE',
    fontSize: 11,
    fontWeight: '600',
  },
});
