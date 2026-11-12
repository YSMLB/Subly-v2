import React from 'react';
import { StyleSheet, View, TouchableOpacity, Text } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { TabType } from '../types';

interface BottomNavBarProps {
  currentTab: TabType;
  onSelectTab: (tab: TabType) => void;
}

export const BottomNavBar: React.FC<BottomNavBarProps> = ({ currentTab, onSelectTab }) => {
  return (
    <View style={styles.wrapper}>
      <View style={styles.container}>
        <TouchableOpacity
          style={[styles.tab, currentTab === 'tasks' && styles.activeTab]}
          onPress={() => onSelectTab('tasks')}
          activeOpacity={0.7}
        >
          <Ionicons
            name={currentTab === 'tasks' ? 'checkbox' : 'checkbox-outline'}
            size={22}
            color={currentTab === 'tasks' ? '#00E5FF' : '#7E8696'}
          />
          <Text style={[styles.tabLabel, currentTab === 'tasks' && styles.activeTabLabel]}>
            Задачи
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[styles.tab, currentTab === 'schedule' && styles.activeTab]}
          onPress={() => onSelectTab('schedule')}
          activeOpacity={0.7}
        >
          <Ionicons
            name={currentTab === 'schedule' ? 'calendar' : 'calendar-outline'}
            size={22}
            color={currentTab === 'schedule' ? '#00E5FF' : '#7E8696'}
          />
          <Text style={[styles.tabLabel, currentTab === 'schedule' && styles.activeTabLabel]}>
            Пары
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[styles.tab, currentTab === 'moderation' && styles.activeTab]}
          onPress={() => onSelectTab('moderation')}
          activeOpacity={0.7}
        >
          <Ionicons
            name={currentTab === 'moderation' ? 'shield-checkmark' : 'shield-checkmark-outline'}
            size={22}
            color={currentTab === 'moderation' ? '#00E5FF' : '#7E8696'}
          />
          <Text style={[styles.tabLabel, currentTab === 'moderation' && styles.activeTabLabel]}>
            Замены
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  wrapper: {
    position: 'absolute',
    bottom: 24,
    left: 20,
    right: 20,
    alignItems: 'center',
  },
  container: {
    flexDirection: 'row',
    backgroundColor: '#121622',
    borderRadius: 36,
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderWidth: 1,
    borderColor: '#202738',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 10 },
    shadowOpacity: 0.6,
    shadowRadius: 16,
    elevation: 12,
    width: '100%',
    justifyContent: 'space-around',
  },
  tab: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 8,
    paddingHorizontal: 16,
    borderRadius: 24,
  },
  activeTab: {
    backgroundColor: '#172233',
    borderWidth: 1,
    borderColor: 'rgba(0, 229, 255, 0.3)',
  },
  tabLabel: {
    color: '#7E8696',
    fontSize: 13,
    fontWeight: '600',
    marginLeft: 6,
  },
  activeTabLabel: {
    color: '#00E5FF',
    fontWeight: '700',
  },
});
