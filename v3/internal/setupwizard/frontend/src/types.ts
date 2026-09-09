export interface DependencyStatus {
  name: string;
  installed: boolean;
  version?: string;
  path?: string;
  status: 'installed' | 'not_installed' | 'needs_update' | 'needs_config' | 'checking';
  required: boolean;
  message?: string;
  installCommand?: string;
  configCommand?: string;
  helpUrl?: string;
  helpLabel?: string;
  imageBuilt?: boolean; // For Docker: whether wails-cross image exists
}

export interface DockerStatus {
  installed: boolean;
  running: boolean;
  version?: string;
  imageBuilt: boolean;
  imageName: string;
  imageSize?: string;
  pullProgress: number;
  pullMessage?: string;
  pullStatus: 'idle' | 'pulling' | 'complete' | 'error';
  pullError?: string;
  bytesTotal?: number;
  bytesDone?: number;
  layerCount?: number;
  layersDone?: number;
}

export interface UserConfig {
  developerName: string;
  email: string;
  defaultFramework: string;
  projectDirectory: string;
  editor: string;
}

export interface WailsConfig {
  info: {
    companyName: string;
    productName: string;
    productIdentifier: string;
    description: string;
    copyright: string;
    comments: string;
    version: string;
  };
}

export interface SystemInfo {
  os: string;
  arch: string;
  wailsVersion: string;
  goVersion: string;
  homeDir: string;
  osName?: string;
  osVersion?: string;
  gitName?: string;
  gitEmail?: string;
}

export interface WizardState {
  currentStep: number;
  dependencies: DependencyStatus[];
  docker: DockerStatus;
  config: UserConfig;
  system: SystemInfo;
  startTime: string;
}

export type Step = 'splash' | 'welcome' | 'dependencies' | 'docker' | 'defaults' | 'signing' | 'config' | 'wails-config' | 'complete';

export interface AuthorDefaults {
  name: string;
  company: string;
}

export interface ProjectDefaults {
  productIdentifierPrefix: string;
  defaultTemplate: string;
  copyrightTemplate: string;
  descriptionTemplate: string;
  defaultVersion: string;
  useInterfaces: boolean;
}

export interface DarwinSigningDefaults {
  identity?: string;
  teamID?: string;
  keychainProfile?: string;
  entitlements?: string;
  p12Path?: string;
  apiKeyPath?: string;
  apiKeyID?: string;
  apiIssuerID?: string;
}

export interface WindowsSigningDefaults {
  certificatePath?: string;
  thumbprint?: string;
  timestampServer?: string;
  cloudProvider?: string;
  cloudKeyID?: string;
}

export interface LinuxSigningDefaults {
  gpgKeyPath?: string;
  gpgKeyID?: string;
  signRole?: string;
}

export interface SigningDefaults {
  darwin?: DarwinSigningDefaults;
  windows?: WindowsSigningDefaults;
  linux?: LinuxSigningDefaults;
}

export interface GlobalDefaults {
  author: AuthorDefaults;
  project: ProjectDefaults;
  signing?: SigningDefaults;
}

export interface GpgKeyInfo {
  keyID: string;
  uid: string;
}

export interface InitTemplate {
  name: string;
  description: string;
}

export interface InitData {
  mode: string;
  projectName: string;
  templateName: string;
  productName: string;
  productCompany: string;
  productIdentifier: string;
  productDescription: string;
  productVersion: string;
  productCopyright: string;
  productComments: string;
  useInterfaces: boolean;
  baseDir: string;
  templates: InitTemplate[];
  defaultTemplate: string;
}

export interface DarwinSigningStatus {
  hasIdentity: boolean;
  identity?: string;
  identities?: string[];
  hasNotarization: boolean;
  teamID?: string;
  configSource?: string;
  rcodesignAvailable: boolean;
}

export interface WindowsSigningStatus {
  hasCertificate: boolean;
  certificateType?: string;
  hasSignTool: boolean;
  timestampServer?: string;
  configSource?: string;
  osslsigncodeAvailable: boolean;
  opensslAvailable: boolean;
}

export interface LinuxSigningStatus {
  hasGpgKey: boolean;
  gpgKeyID?: string;
  configSource?: string;
  gpgAvailable: boolean;
  gpgKeys?: GpgKeyInfo[];
}

export interface SigningStatus {
  host: 'darwin' | 'windows' | 'linux';
  darwin: DarwinSigningStatus;
  windows: WindowsSigningStatus;
  linux: LinuxSigningStatus;
}
