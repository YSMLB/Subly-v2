import React, { useState } from 'react';
import {
  StyleSheet,
  Text,
  View,
  TextInput,
  TouchableOpacity,
  Image,
  KeyboardAvoidingView,
  Platform,
  TouchableWithoutFeedback,
  Keyboard,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as Haptics from 'expo-haptics';

interface LoginScreenProps {
  onLoginSuccess: () => void;
}

export const LoginScreen: React.FC<LoginScreenProps> = ({ onLoginSuccess }) => {
  // Режим экрана: вход или регистрация
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [registerStep, setRegisterStep] = useState<'form' | 'code'>('form');

  // Поля входа
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  // Поля регистрации
  const [regName, setRegName] = useState('');
  const [regGroup, setRegGroup] = useState('ИТ-31');
  const [regEmail, setRegEmail] = useState('');
  const [regPassword, setRegPassword] = useState('');
  const [showRegPassword, setShowRegPassword] = useState(false);
  const [verifyCode, setVerifyCode] = useState('');

  // Состояние загрузки
  const [loading, setLoading] = useState(false);

  const handleSignIn = async () => {
    Keyboard.dismiss();
    await Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Medium);
    onLoginSuccess();
  };

  const handleSendCode = async () => {
    if (!regEmail.trim()) {
      Alert.alert('Ошибка', 'Введите email для отправки кода подтверждения.');
      return;
    }
    if (!regName.trim()) {
      Alert.alert('Ошибка', 'Введите имя или логин.');
      return;
    }

    Keyboard.dismiss();
    setLoading(true);
    await Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);

    // Имитация отправки запроса на бэкенд
    setTimeout(async () => {
      setLoading(false);
      setRegisterStep('code');
      await Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
      Alert.alert(
        'Код отправлен!',
        `6-значный код подтверждения отправлен на ${regEmail}.\n(Для быстрого теста: 742195)`,
        [{ text: 'Ввести код' }]
      );
    }, 1000);
  };

  const handleVerifyAndRegister = async () => {
    if (!verifyCode.trim()) {
      Alert.alert('Ошибка', 'Введите 6-значный код из письма.');
      return;
    }

    setLoading(true);
    await Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Medium);

    setTimeout(async () => {
      setLoading(false);
      await Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
      Alert.alert('Успешно', 'Аккаунт успешно создан! Добро пожаловать в Subly.', [
        {
          text: 'Войти в приложение',
          onPress: onLoginSuccess,
        },
      ]);
    }, 800);
  };

  return (
    <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.container}
      >
        <View style={styles.header}>
          <Image
            source={require('../../assets/subly-logo.jpg')}
            style={styles.logo}
            resizeMode="contain"
          />
          <Text style={styles.title}>Subly</Text>
          <Text style={styles.subtitle}>
            {mode === 'login' ? 'Твой день под контролем' : 'Создание нового аккаунта'}
          </Text>
        </View>

        {mode === 'login' ? (
          /* =================== ФОРМА ВХОДА =================== */
          <View style={styles.card}>
            <View style={styles.inputContainer}>
              <TextInput
                style={styles.input}
                placeholder="Логин или телефон"
                placeholderTextColor="#717A8C"
                value={login}
                onChangeText={setLogin}
                autoCapitalize="none"
              />
            </View>

            <View style={[styles.inputContainer, styles.passwordContainer]}>
              <TextInput
                style={styles.input}
                placeholder="Пароль"
                placeholderTextColor="#717A8C"
                secureTextEntry={!showPassword}
                value={password}
                onChangeText={setPassword}
              />
              <TouchableOpacity
                onPress={() => setShowPassword(!showPassword)}
                hitSlop={{ top: 12, bottom: 12, left: 12, right: 12 }}
              >
                <Ionicons
                  name={showPassword ? 'eye-off-outline' : 'eye-outline'}
                  size={21}
                  color="#8A93A6"
                />
              </TouchableOpacity>
            </View>

            <TouchableOpacity
              style={styles.button}
              activeOpacity={0.8}
              onPress={handleSignIn}
            >
              <Text style={styles.buttonText}>Войти</Text>
            </TouchableOpacity>

            <View style={styles.footer}>
              <TouchableOpacity
                onPress={() =>
                  Alert.alert(
                    'Сброс пароля',
                    'Для восстановления доступа напишите разработчику в Telegram: @amiri2sb'
                  )
                }
              >
                <Text style={styles.footerLink}>Забыли пароль?</Text>
              </TouchableOpacity>

              <TouchableOpacity
                onPress={() => {
                  setMode('register');
                  setRegisterStep('form');
                }}
              >
                <Text style={[styles.footerLink, styles.footerHighlight]}>
                  Создать аккаунт
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        ) : (
          /* =================== ФОРМА РЕГИСТРАЦИИ =================== */
          <View style={styles.card}>
            {registerStep === 'form' ? (
              <>
                <View style={styles.inputContainer}>
                  <TextInput
                    style={styles.input}
                    placeholder="Имя или логин"
                    placeholderTextColor="#717A8C"
                    value={regName}
                    onChangeText={setRegName}
                    autoCapitalize="none"
                  />
                </View>

                <View style={styles.inputContainer}>
                  <TextInput
                    style={styles.input}
                    placeholder="Учебная группа (например, ИТ-31)"
                    placeholderTextColor="#717A8C"
                    value={regGroup}
                    onChangeText={setRegGroup}
                    autoCapitalize="characters"
                  />
                </View>

                <View style={styles.inputContainer}>
                  <TextInput
                    style={styles.input}
                    placeholder="Email (для кода подтверждения)"
                    placeholderTextColor="#717A8C"
                    value={regEmail}
                    onChangeText={setRegEmail}
                    keyboardType="email-address"
                    autoCapitalize="none"
                  />
                </View>

                <View style={[styles.inputContainer, styles.passwordContainer]}>
                  <TextInput
                    style={styles.input}
                    placeholder="Придумайте пароль"
                    placeholderTextColor="#717A8C"
                    secureTextEntry={!showRegPassword}
                    value={regPassword}
                    onChangeText={setRegPassword}
                  />
                  <TouchableOpacity
                    onPress={() => setShowRegPassword(!showRegPassword)}
                    hitSlop={{ top: 12, bottom: 12, left: 12, right: 12 }}
                  >
                    <Ionicons
                      name={showRegPassword ? 'eye-off-outline' : 'eye-outline'}
                      size={21}
                      color="#8A93A6"
                    />
                  </TouchableOpacity>
                </View>

                <TouchableOpacity
                  style={styles.button}
                  activeOpacity={0.8}
                  disabled={loading}
                  onPress={handleSendCode}
                >
                  {loading ? (
                    <ActivityIndicator size="small" color="#000000" />
                  ) : (
                    <Text style={styles.buttonText}>Получить код на почту</Text>
                  )}
                </TouchableOpacity>
              </>
            ) : (
              /* Ввод кода подтверждения */
              <>
                <View style={styles.codePromptBox}>
                  <Ionicons name="mail-unread-outline" size={28} color="#00E5FF" />
                  <Text style={styles.codePromptTitle}>Подтверждение Email</Text>
                  <Text style={styles.codePromptText}>
                    Мы отправили 6-значный код на{'\n'}
                    <Text style={styles.codePromptEmail}>{regEmail}</Text>
                  </Text>
                </View>

                <View style={styles.inputContainer}>
                  <TextInput
                    style={[styles.input, styles.codeInput]}
                    placeholder="• • • • • •"
                    placeholderTextColor="#5E6778"
                    value={verifyCode}
                    onChangeText={setVerifyCode}
                    keyboardType="number-pad"
                    maxLength={6}
                    autoFocus
                  />
                </View>

                <TouchableOpacity
                  style={styles.button}
                  activeOpacity={0.8}
                  disabled={loading}
                  onPress={handleVerifyAndRegister}
                >
                  {loading ? (
                    <ActivityIndicator size="small" color="#000000" />
                  ) : (
                    <Text style={styles.buttonText}>Завершить регистрацию</Text>
                  )}
                </TouchableOpacity>

                <TouchableOpacity
                  style={styles.resendBtn}
                  onPress={() => setRegisterStep('form')}
                >
                  <Text style={styles.resendText}>← Изменить email или данные</Text>
                </TouchableOpacity>
              </>
            )}

            <View style={styles.footerCenter}>
              <TouchableOpacity onPress={() => setMode('login')}>
                <Text style={styles.footerLink}>
                  Уже есть аккаунт? <Text style={styles.footerHighlight}>Войти</Text>
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        )}
      </KeyboardAvoidingView>
    </TouchableWithoutFeedback>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000000',
    justifyContent: 'center',
    paddingHorizontal: 24,
  },
  header: {
    alignItems: 'center',
    marginBottom: 28,
  },
  logo: {
    width: 90,
    height: 90,
    marginBottom: 12,
  },
  title: {
    fontSize: 40,
    fontWeight: '800',
    color: '#FFFFFF',
    letterSpacing: -0.5,
  },
  subtitle: {
    fontSize: 14,
    color: '#8A93A6',
    marginTop: 4,
  },
  card: {
    backgroundColor: '#121622',
    borderRadius: 24,
    padding: 22,
    borderWidth: 1,
    borderColor: '#202738',
    shadowColor: '#000000',
    shadowOffset: { width: 0, height: 10 },
    shadowOpacity: 0.5,
    shadowRadius: 20,
    elevation: 8,
  },
  inputContainer: {
    backgroundColor: '#171D2B',
    borderRadius: 14,
    height: 52,
    marginBottom: 14,
    paddingHorizontal: 16,
    justifyContent: 'center',
    borderWidth: 1,
    borderColor: '#263145',
  },
  passwordContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  input: {
    color: '#FFFFFF',
    fontSize: 15,
    flex: 1,
  },
  codeInput: {
    textAlign: 'center',
    fontSize: 22,
    letterSpacing: 8,
    fontWeight: '700',
    color: '#00E5FF',
  },
  button: {
    backgroundColor: '#00E5FF',
    height: 52,
    borderRadius: 14,
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 8,
    marginBottom: 18,
    shadowColor: '#00E5FF',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.6,
    shadowRadius: 14,
    elevation: 6,
  },
  buttonText: {
    color: '#000000',
    fontSize: 16,
    fontWeight: '800',
    letterSpacing: 0.5,
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingHorizontal: 4,
  },
  footerCenter: {
    alignItems: 'center',
    paddingTop: 4,
  },
  footerLink: {
    color: '#7E8696',
    fontSize: 13,
    fontWeight: '500',
  },
  footerHighlight: {
    color: '#00E5FF',
    fontWeight: '700',
  },
  codePromptBox: {
    alignItems: 'center',
    marginBottom: 18,
    paddingVertical: 8,
  },
  codePromptTitle: {
    color: '#FFFFFF',
    fontSize: 17,
    fontWeight: '700',
    marginTop: 8,
  },
  codePromptText: {
    color: '#7E8696',
    fontSize: 13,
    textAlign: 'center',
    marginTop: 4,
    lineHeight: 18,
  },
  codePromptEmail: {
    color: '#00E5FF',
    fontWeight: '600',
  },
  resendBtn: {
    alignItems: 'center',
    marginBottom: 14,
  },
  resendText: {
    color: '#7E8696',
    fontSize: 12,
  },
});
