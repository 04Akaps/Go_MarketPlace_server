# MarketPlace를 만드는 서버

특정 서명에 대한 복호화, 채팅 기능 및 추가적인 기능은 작성하면서 사용할 예정

DB는 두개 사용할 예정
- mysql, mongo -> dynamoDB로 변경도 가능

캐시는 Redis사용 예정

Jenkins로 배포하고, 이번에는 Docker가 아니라 EC2의 로드 밸런싱 활용 예정
- CloudFront도 가능하다면??

이런식이면 실제로 비용이 부과될 것이니.. 일단 최대한 로컬망으로 작업 후 이후 EC2연동
- 사양한 가능하다면 프리티어로;;



# 사용하는 패키지
Viper
Gomniauth
Paseto
...++


# 버전 관리
glide

새로운 패키지 추가 하는 법
1. glide get <패키지>
2. version track은 Patch로 사용하여 최소화
3. glide install

패키지 업데이트 할 경우 사용하기
1. glide up

# env 관리

내부에 app.env로 설정하여 다음과 같은 설정을 실행

