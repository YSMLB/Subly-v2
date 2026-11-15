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

export default function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [currentTab, setCurrentTab] = useState<TabType>('tasks');
  const [bannerNotification, setBannerNotification] = useState<NotificationPayload | null>(null);

  useEffect(() => {
    requestNotificationPermissions();
  }, []);

  const handleTriggerNotification = (title: string, message: string, type: 'replacement' | 'reminder' | 'system' | 'task' = 'task') => {
    setBannerNotification({
      id: Date.now().toString(),
      title,
      message,
      type,
    });
    sendPushNotification(title, message);
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

