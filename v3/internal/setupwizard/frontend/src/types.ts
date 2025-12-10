export interface DependencyStatus {
  name: string;
  installed: boolean;
  version?: string;
  path?: string;
  status: 'installed' | 'not_installed' | 'needs_update' | 'checking';
  required: boolean;
  message?: string;
  installCommand?: string;
  helpUrl?: string;
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
  pullStatus: 'idle' | 'pulling' | 'complete' | 'error';
  pullError?: string;
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
}

export interface SigningDefaults {
  macOS: {
    developerID: string;        // e.g., "Developer ID Application: John Doe (TEAMID)"
    appleID: string;            // Apple ID for notarization
    teamID: string;             // Apple Team ID
  };
  windows: {
    certificatePath: string;    // Path to .pfx certificate
    timestampServer: string;    // e.g., "http://timestamp.digicert.com"
  };
}

export interface GlobalDefaults {
  author: AuthorDefaults;
  project: ProjectDefaults;
  signing?: SigningDefaults;
}
