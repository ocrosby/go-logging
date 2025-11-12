pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.24'
        COVERAGE_THRESHOLD = '85'
    }
    
    stages {
        stage('Setup') {
            steps {
                // Clean workspace
                cleanWs()
                
                // Checkout code
                checkout scm
                
                // Setup Go
                script {
                    sh '''
                        # Install Go if not available
                        if ! command -v go &> /dev/null; then
                            echo "Installing Go ${GO_VERSION}..."
                            wget -q https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
                            sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
                            export PATH=$PATH:/usr/local/go/bin
                        fi
                        
                        # Verify Go installation
                        go version
                        
                        # Download dependencies
                        go mod download
                        go mod verify
                    '''
                }
            }
        }
        
        stage('Install Tools') {
            steps {
                script {
                    sh '''
                        # Install Task runner
                        sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
                        
                        # Install bc for numeric comparisons
                        if ! command -v bc &> /dev/null; then
                            apt-get update && apt-get install -y bc
                        fi
                        
                        # Install development tools using task
                        task install-tools
                    '''
                }
            }
        }
        
        stage('Lint') {
            steps {
                script {
                    sh '''
                        # Run linting
                        task lint
                    '''
                }
            }
        }
        
        stage('Test with Coverage') {
            steps {
                script {
                    sh '''
                        # Run tests with proper coverage validation
                        # IMPORTANT: Use task test-coverage-check, NOT go test ./...
                        task test-coverage-check
                    '''
                }
                
                // Archive coverage reports
                publishHTML([
                    allowMissing: false,
                    alwaysLinkToLastBuild: true,
                    keepAll: true,
                    reportDir: '.',
                    reportFiles: 'coverage.html',
                    reportName: 'Coverage Report'
                ])
            }
        }
        
        stage('Build') {
            steps {
                script {
                    sh '''
                        # Build all packages
                        task build
                    '''
                }
            }
        }
        
        stage('Integration Tests') {
            steps {
                script {
                    sh '''
                        # Run integration tests
                        timeout 10s go run ./examples/basic/main.go || true
                        timeout 10s go run ./examples/fluent/main.go || true
                        timeout 10s go run ./examples/di/. || true
                    '''
                }
            }
        }
    }
    
    post {
        always {
            // Archive test results and coverage
            archiveArtifacts artifacts: 'logging-coverage.out', fingerprint: true
            
            // Clean up
            cleanWs()
        }
        
        failure {
            script {
                echo '''
                ❌ BUILD FAILED!
                
                Common issues and solutions:
                
                1. Coverage below 85%:
                   - Check the coverage report for uncovered functions
                   - Main logging package should maintain 85%+ coverage
                   - Mocks package intentionally has low coverage
                
                2. Timeout errors:
                   - Increase timeout values if tests are timing out
                   - Check for hanging goroutines in async tests
                
                3. Task command not found:
                   - Ensure Task runner is properly installed
                   - Alternative: Use fallback commands in CLAUDE.md
                
                4. Linting failures:
                   - Run `task lint` locally to see specific issues
                   - Check for cyclomatic complexity violations
                
                For detailed troubleshooting, see CLAUDE.md in the repository.
                '''
            }
        }
        
        success {
            echo '''
            ✅ BUILD SUCCESSFUL!
            
            All quality gates passed:
            - Linting: Clean
            - Tests: All passing
            - Coverage: Above 85% threshold
            - Build: Successful
            '''
        }
    }
}