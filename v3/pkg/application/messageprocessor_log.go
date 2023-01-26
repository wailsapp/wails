package application

//
//import "errors"
//
////var logLevelMap = map[byte]logger.LogLevel{
////	'1': pkgLogger.TRACE,
////	'2': pkgLogger.DEBUG,
////	'3': pkgLogger.INFO,
////	'4': pkgLogger.WARNING,
////	'5': pkgLogger.ERROR,
////}
//
//func (m *MessageProcessor) processLogMessage(message string) {
//	if len(message) < 3 {
//		m.Error("Invalid Log Message: " + message)
//		return
//	}
//
//	messageText := message[2:]
//
//	switch message[1] {
//	case 'T':
//		d.log.Trace(messageText)
//	case 'P':
//		d.log.Print(messageText)
//	case 'D':
//		d.log.Debug(messageText)
//	case 'I':
//		d.log.Info(messageText)
//	case 'W':
//		d.log.Warning(messageText)
//	case 'E':
//		d.log.Error(messageText)
//	case 'F':
//		d.log.Fatal(messageText)
//	case 'S':
//		loglevel, exists := logLevelMap[message[2]]
//		if !exists {
//			return "", errors.New("Invalid Set Log Level Message: " + message)
//		}
//		d.log.SetLogLevel(loglevel)
//	default:
//		return "", errors.New("Invalid Log Message: " + message)
//	}
//	return "", nil
//}
