import React, { useState, useEffect } from 'react';
import { StyleSheet, SafeAreaView, View } from 'react-native';
import { StatusBar } from 'expo-status-bar';
import { TabType, NotificationPayload } from './src/types';
import { LoginScreen } from './src/screens/LoginScreen';
import { TasksScreen } from './src/screens/TasksScreen';
import { ScheduleScreen } from './src/screens/ScheduleScreen';
import { ModerationScreen } from './src/screens/ModerationScreen';
import { BottomNavBar } from './src/components/BottomNavBar';
import { NotificationBanner } from './src/components/NotificationBanner';
import { requestNotificationPermissions, sendPushNotification } from './src/services/notificationService';
import { getAuthSession, clearAuthSession } from './src/api/client';

export default function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [currentTab, setCurrentTab] = useState<TabType>('tasks');
  const [bannerNotification, setBannerNotification] = useState<NotificationPayload | null>(null);

  useEffect(() => {
    requestNotificationPermissions();
    getAuthSession().then(({ token }) => {
      if (token) {
        setIsLoggedIn(true);
      }
    });
  }, []);

  const handleLogout = async () => {
    await clearAuthSession();
    setIsLoggedIn(false);
  };

  const handleTriggerNotification = (
    title: string,
    message: string,
    type: 'replacement' | 'reminder' | 'system' | 'task' = 'task'
  ) => {
    // Внутренний неоновый баннер показывает полные данные со свечением:
    setBannerNotification({
      id: Date.now().toString(),
      title,
      message,
      type,
    });
    // Системный пуш на экран блокировки в стелс-режиме (без палева):
    sendPushNotification('Subly', 'Есть изменения', true);
  };

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar style="light" />

      <NotificationBanner
        notification={bannerNotification}
        onDismiss={() => setBannerNotification(null)}
      />

      {!isLoggedIn ? (
        <LoginScreen onLoginSuccess={() => setIsLoggedIn(true)} />
      ) : (
        <View style={styles.mainContainer}>
          {currentTab === 'tasks' && (
            <TasksScreen
              onNavigateToSchedule={() => setCurrentTab('schedule')}
              onTriggerNotification={handleTriggerNotification}
              onLogout={handleLogout}
            />
          )}
          {currentTab === 'schedule' && <ScheduleScreen />}
          {currentTab === 'moderation' && (
            <ModerationScreen
              onTriggerNotification={handleTriggerNotification}
            />
          )}

          <BottomNavBar currentTab={currentTab} onSelectTab={setCurrentTab} />
        </View>
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000000',
  },
  mainContainer: {
    flex: 1,
    position: 'relative',
    backgroundColor: '#000000',
  },
});

