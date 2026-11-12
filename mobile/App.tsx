import React, { useState } from 'react';
import { StyleSheet, SafeAreaView, View } from 'react-native';
import { StatusBar } from 'expo-status-bar';
import { TabType } from './src/types';
import { LoginScreen } from './src/screens/LoginScreen';
import { TasksScreen } from './src/screens/TasksScreen';
import { ScheduleScreen } from './src/screens/ScheduleScreen';
import { ModerationScreen } from './src/screens/ModerationScreen';
import { BottomNavBar } from './src/components/BottomNavBar';

export default function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [currentTab, setCurrentTab] = useState<TabType>('tasks');

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar style="light" />
      {!isLoggedIn ? (
        <LoginScreen onLoginSuccess={() => setIsLoggedIn(true)} />
      ) : (
        <View style={styles.mainContainer}>
          {currentTab === 'tasks' && <TasksScreen />}
          {currentTab === 'schedule' && <ScheduleScreen />}
          {currentTab === 'moderation' && <ModerationScreen />}

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
